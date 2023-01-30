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
	rpsChan := deployRpsComposer(ctx, startTime, &csp.Config)

	go deployDefaultLoadsRequester(ctx, cancelCtx, defaultRequesterChan)
	go deployCustomLoadsRequester(ctx, cancelCtx, customRequesterChan)

	currentRps := int(getInitialRps(&csp.Config))
	previousRps := currentRps
	defaultRequesterRps := currentRps
l1:
	for {
		now := time.Now()
		runtime := time.Now().Unix() - startTime.Unix()

		select {
		case <-ctx.Done():
			common.GetLogger().Log("Execution finished")
			common.GetLogger().RegisterLogs()
			break l1
		case newRps := <-rpsChan:
			customLoad := isCustomLoadWindow(&csp.Config, now.Unix())
			if newRps != 0 {
				previousRps = currentRps
				currentRps = newRps
				defaultRequesterRps = currentRps

				if customLoad == nil {
					logRps(previousRps, currentRps, runtime)

					defaultRequesterChan <- DefaultRequesterPayload{
						Request: requestQueue[requestQueueIter.Next()],
						Rps:     currentRps,
					}
				}
			}
		default:
			customLoad := isCustomLoadWindow(&csp.Config, now.Unix())
			request := requestQueue[requestQueueIter.Next()]

			if customLoad != nil {
				previousRps = currentRps
				currentRps = customLoad.Rps

				if previousRps != currentRps {
					common.GetLogger().Log(fmt.Sprintf("Runtime: %ds, Rps: %d (CUSTOM)", runtime, customLoad.Rps))
				}

				customRequesterChan <- CustomRequesterPayload{
					Request:          request,
					CustomLoadConfig: customLoad,
				}
			} else {
				previousRps = currentRps
				currentRps = defaultRequesterRps

				logRps(previousRps, currentRps, runtime)
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

func logRps(prevRps, currentRps int, runtime int64) {
	if prevRps != currentRps {
		common.GetLogger().Log(fmt.Sprintf("Runtime: %ds, Rps: %d", runtime, currentRps))
	}
}

// This one purpose is to keep track of whenever a rampup should be made
func deployRpsComposer(ctx context.Context, startTime time.Time, cpc *CustomProfileConfig) <-chan int {
	rpsChan := make(chan int)
	ticker := time.NewTicker(cpc.RpsIncreaseInterval.Duration)

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

					common.GetLogger().Log("Rampup finished")
					wg.Done()
					break l1
				} else {
					rawRps += cpc.RpsRampupPace
					effectiveRps = int(rawRps)

					if effectiveRps != 0 {
						rpsChan <- effectiveRps
					}
				}
			}

		}
	}()

	return rpsChan
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
	for i, v := range csp.Config.CustomLoads {
		v.calculateWindow(startTime)

		csp.Config.CustomLoads[i] = v
	}
}

func getInitialRps(cpc *CustomProfileConfig) float64 {
	tempRps := cpc.RpsRampupPace

	if tempRps > 1 {
		return tempRps
	}

	return 1
}
