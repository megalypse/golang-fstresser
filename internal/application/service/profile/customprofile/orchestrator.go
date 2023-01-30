package customprofile

import (
	"context"
	"fmt"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/common"
)

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
	requestCountChan := make(chan int)
	rpsChan := deployRpsComposer(ctx, startTime, &csp.Config)

	wg.Add(1)
	go deployErrorThresholdAnalyzer(ctx, cancelCtx, requestCountChan, &csp.Config)

	wg.Add(1)
	go deployDefaultRequester(ctx, cancelCtx, defaultRequesterChan, requestCountChan)

	wg.Add(1)
	go deployCustomRequester(ctx, cancelCtx, customRequesterChan, requestCountChan)

	currentRps := int(getInitialRps(&csp.Config))
	previousRps := currentRps
	defaultRequesterRps := currentRps
l1:
	for {
		now := time.Now()
		runtime := time.Now().Unix() - startTime.Unix()

		select {
		case <-ctx.Done():
			common.GetLogger().Log(fmt.Sprintf("Execution finished. Took %ds", runtime))
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
			if now.Unix() >= startTime.Add(csp.Config.EndLoadAt.Duration).Unix() {
				cancelCtx()
				continue
			} else {
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
		}

		time.Sleep(time.Second)
	}

	wg.Wait()

	close(defaultRequesterChan)
	close(customRequesterChan)
	close(requestCountChan)
}
