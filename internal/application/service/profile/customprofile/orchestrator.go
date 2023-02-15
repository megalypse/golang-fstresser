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
	defer common.HandlePanic(ctx, cancelCtx)

	startTime := time.Now()

	csp.Config.bootstrap()
	calculateCustomLoadsWindows(startTime, csp)
	prepareRequests(csp)

	requestQueueIter := generateRequestQueue(len(requestQueue) - 1)

	defaultRequesterChan := make(chan DefaultRequesterPayload)
	customRequesterChan := make(chan CustomRequesterPayload)
	requestCountChan := make(chan int)

	wg.Add(1)
	go deployHttpStatusAnalyzer(ctx, cancelCtx, requestCountChan, &csp.Config)

	wg.Add(1)
	go deployDefaultRequester(ctx, cancelCtx, csp, defaultRequesterChan, requestCountChan)

	wg.Add(1)
	go deployCustomRequester(ctx, cancelCtx, csp, customRequesterChan, requestCountChan)

	rpsChan := deployRpsComposer(ctx, cancelCtx, startTime, &csp.Config)

	currentRps := int(getInitialRps(&csp.Config))
	previousRps := currentRps
	defaultRequesterRps := currentRps
l1:
	for {
		now := time.Now()
		runtime := time.Now().Unix() - startTime.Unix()
		durationRuntime := time.Second * time.Duration(runtime)

		select {
		case <-ctx.Done():
			common.GetLogger(ctx).Log(fmt.Sprintf("(%s) Execution finished. Took %ds", ctx.Value(common.GetCtxKey("profile-name")), runtime))
			common.GetLogger(ctx).RegisterLogs()
			break l1
		case newRps := <-rpsChan:
			customLoad := isCustomLoadWindow(&csp.Config, now.Unix())
			if newRps != 0 {
				previousRps = currentRps
				currentRps = newRps
				defaultRequesterRps = currentRps

				if customLoad == nil {
					logRps(ctx, previousRps, currentRps, durationRuntime)

					defaultRequesterChan <- DefaultRequesterPayload{
						Request: requestQueue[requestQueueIter.Next()],
						Rps:     currentRps,
					}
				}
			}
		default:
			if now.Unix() >= startTime.Add(csp.Config.ExecutionEndsAt.Duration).Unix() {
				cancelCtx()
				continue
			} else {
				customLoad := isCustomLoadWindow(&csp.Config, now.Unix())
				request := requestQueue[requestQueueIter.Next()]

				if customLoad != nil {
					previousRps = currentRps
					currentRps = customLoad.Rps

					if previousRps != currentRps {
						common.GetLogger(ctx).Log(fmt.Sprintf("(%s) Runtime: %s, Rps: %d (CUSTOM)", ctx.Value(common.GetCtxKey("profile-name")), durationRuntime.String(), currentRps))
					}

					customRequesterChan <- CustomRequesterPayload{
						Request:          request,
						CustomLoadConfig: customLoad,
					}
				} else {
					previousRps = currentRps
					currentRps = defaultRequesterRps

					logRps(ctx, previousRps, currentRps, durationRuntime)
					defaultRequesterChan <- DefaultRequesterPayload{
						Request: request,
						Rps:     currentRps,
					}
				}
			}
		}

		time.Sleep(csp.Config.LoadsInterval.Duration)
	}

	wg.Wait()

	defer close(defaultRequesterChan)
	defer close(customRequesterChan)
	defer close(requestCountChan)
	defer close(rpsChan)
}
