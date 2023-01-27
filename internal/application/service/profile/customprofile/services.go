package customprofile

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

// This mutex is created to be used inside this file only, please,
// DO NOT USE IT ANYWHERE ELSE
var wg sync.WaitGroup
var requestQueue []*entity.Request

func init() {
	requestQueue = make([]*entity.Request, 0)
}

func deployCustomProfileOrchestrator(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	csp *CustomStressProfile,
) {
	startTime := time.Now()

	csp.Config.bootstrap()
	calculateCustomLoadsWindows(startTime, csp)
	prepareRequests(csp)
	requestQueueIter := generateRequestQueue(len(requestQueue) - 1)

	defaultRequesterChan := make(chan DefaultRequesterPayload)
	customRequesterChan := make(chan CustomRequesterPayload)
	rpsChan, _ := deployRpsComposer(ctx, startTime, &csp.Config)

	go deployDefaultLoadsRequester(ctx, cancelCtx, defaultRequesterChan)
	go deployCustomLoadsRequester(ctx, cancelCtx, customRequesterChan)

	currentRps := int(getInitialRps(&csp.Config))
	previousRps := currentRps
l1:
	for {
		now := time.Now()
		customLoad := isCustomLoadWindow(&csp.Config, now.Unix())

		select {
		case <-ctx.Done():
			common.GetLogger().Log("Execution finished")
			common.GetLogger().RegisterLogs()
			break l1
		case newRps := <-rpsChan:
			previousRps = currentRps
			currentRps = newRps

			if customLoad == nil {
				logRps(previousRps, currentRps, fmt.Sprintf("Rps: %d", currentRps))

				defaultRequesterChan <- DefaultRequesterPayload{
					Request: requestQueue[requestQueueIter.Next()],
					Rps:     currentRps,
				}
			}
		default:
			request := requestQueue[requestQueueIter.Next()]

			if customLoad != nil {
				common.GetLogger().Log(fmt.Sprintf("Rps: %d (CUSTOM)", customLoad.Rps))

				customRequesterChan <- CustomRequesterPayload{
					Request:          request,
					CustomLoadConfig: customLoad,
				}
			} else {
				logRps(previousRps, currentRps, fmt.Sprintf("Rps: %d", currentRps))

				defaultRequesterChan <- DefaultRequesterPayload{
					Request: request,
					Rps:     currentRps,
				}
			}
		}

		time.Sleep(time.Second)
	}

	wg.Wait()

	close(defaultRequesterChan)
	close(customRequesterChan)
}

func logRps(prevRps, currentRps int, message string) {
	if prevRps != currentRps {
		common.GetLogger().Log(message)
	}
}

// This one purpose is to keep track of whenever a rampup should be made
func deployRpsComposer(ctx context.Context, startTime time.Time, cpc *CustomProfileConfig) (<-chan int, context.Context) {
	rpsChan := make(chan int)
	ticker := time.NewTicker(cpc.RpsIncreaseInterval.Duration)
	tickerCtx, cancelTickerCtx := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		rawRps := getInitialRps(cpc)
		effectiveRps := int(rawRps)
		tickerDeadline := startTime.Add(cpc.RampUpTime.Duration).Unix()
	l1:
		for {
			select {
			case <-ctx.Done():
				close(rpsChan)

				wg.Done()
				break l1
			case <-ticker.C:
				// In case the desired peak RPS have not been met at the end of the rampup time
				nowUnix := time.Now().Unix()
				if nowUnix > tickerDeadline {
					if effectiveRps != cpc.PeakRps {
						rpsChan <- cpc.PeakRps
					}

					close(rpsChan)
					cancelTickerCtx()
					wg.Done()
					break l1
				} else {
					rawRps += cpc.RpsRampupPace
					effectiveRps = int(rawRps)

					rpsChan <- effectiveRps
				}
			}

		}
	}()

	return rpsChan, tickerCtx
}

func deployDefaultLoadsRequester(ctx context.Context, cancelCtx context.CancelFunc, loadsConsumer <-chan DefaultRequesterPayload) {
	wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case load := <-loadsConsumer:
			for i := 0; i < load.Rps; i++ {
				go common.MakeLightweightRequest(cancelCtx, load.Request)
			}
		}
	}
}

func deployCustomLoadsRequester(ctx context.Context, cancelCtx context.CancelFunc, loadsConsumer <-chan CustomRequesterPayload) {
	wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case load := <-loadsConsumer:
			for i := 0; i < load.CustomLoadConfig.Rps; i++ {
				go common.MakeLightweightRequest(cancelCtx, load.Request)
			}
		}
	}
}

func isCustomLoadWindow(cpc *CustomProfileConfig, now int64) *CustomLoad {
	for _, v := range cpc.CustomLoads {
		if now >= v.GetStartPoint() && now < v.GetEndPoint() {
			return &v
		}
	}

	return nil
}

func generateRequestQueue(maxCount int) RateCounter {
	if maxCount < 0 {
		common.GetLogger().Log("Max count must be a positive number")
	}

	return RateCounter{
		maxValue:     maxCount,
		currentValue: -1,
	}
}

func prepareRequests(csp *CustomStressProfile) {
	for _, v := range csp.Requests {
		requestEntity := v.ToEntity()

		for i := 0; i < v.Rate; i++ {
			requestQueue = append(requestQueue, &requestEntity)
		}
	}
}

func calculateCustomLoadsWindows(startTime time.Time, csp *CustomStressProfile) {
	for _, v := range csp.Config.CustomLoads {
		v.calculateWindow(startTime)
	}
}

func getInitialRps(cpc *CustomProfileConfig) float64 {
	tempRps := cpc.RpsRampupPace

	if tempRps > 1 {
		return tempRps
	}

	return 1
}
