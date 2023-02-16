package customprofile

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/megalypse/golang-fstresser/internal/application/common"
)

/*
deployHttpStatusAnalyzer should be executed on its own routine.

Its goal is to analyze the ratio of succeeded/failed requests,
and gracefully shutdown the software execution when the threshold is met.
*/
func deployHttpStatusAnalyzer(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	reqCountConsumer <-chan int,
	cpc *CustomProfileConfig,
) {
	defer wg.Done()
	defer common.HandlePanic(ctx, cancelCtx)

	var successfullRequests int
	var failedRequests int
	var totalRequests int

	separated := strings.Split(cpc.ErrorThreshold, ":")
	threshold, err := strconv.Atoi(separated[0])
	if err != nil {
		common.GetLogger(ctx).Log(err.Error())
		cancelCtx()
		return
	}

	for {
		select {
		case <-ctx.Done():
			message := fmt.Sprintf("(%s) Requests done: %d\nSuccess: %d\nFailed: %d", ctx.Value(common.GetCtxKey("profile-name")), totalRequests, successfullRequests, failedRequests)
			common.GetLogger(ctx).Log(message)
			return
		case httpStatus := <-reqCountConsumer:
			if httpStatus >= 200 && httpStatus < 300 {
				successfullRequests++
			} else {
				failedRequests++
			}

			totalRequests = successfullRequests + failedRequests
			errMsg := fmt.Sprintf("(%s) Error threshold met", ctx.Value(common.GetCtxKey("profile-name")))

			// The minimum amount of requests for the analysis to happen was set to 10
			// to give at least a small window for it to recover and the software does not
			// shutdown on the first fail.
			if totalRequests >= 10 {
				switch separated[1] {
				case "raw":
					if failedRequests > threshold {
						common.GetLogger(ctx).Log(errMsg)
						cancelCtx()
					}
				default:
					errPercent := failedRequests * 100 / totalRequests

					if errPercent > threshold {
						common.GetLogger(ctx).Log(errMsg)
						cancelCtx()
					}
				}
			}

		}
	}
}
