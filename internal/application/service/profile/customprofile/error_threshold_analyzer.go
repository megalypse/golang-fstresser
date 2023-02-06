package customprofile

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/megalypse/golang-fstresser/internal/application/common"
)

func deployHttpStatusAnalyzer(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	reqCountConsumer <-chan int,
	cpc *CustomProfileConfig,
) {
	defer wg.Done()

	var successfullRequests int
	var failedRequests int
	var totalRequests int

	separated := strings.Split(cpc.ErrorThreshold, ":")
	threshold, err := strconv.Atoi(separated[0])
	if err != nil {
		common.GetLogger().Log(err.Error())
		cancelCtx()
		return
	}

	for {
		select {
		case <-ctx.Done():
			message := fmt.Sprintf("(%s) Requests done: %d\nSuccess: %d\nFailed: %d", ctx.Value("profile-name"), totalRequests, successfullRequests, failedRequests)
			common.GetLogger().Log(message)
			return
		case httpStatus := <-reqCountConsumer:
			if httpStatus >= 200 && httpStatus < 300 {
				successfullRequests++
			} else {
				failedRequests++
			}

			totalRequests = successfullRequests + failedRequests
			errMsg := fmt.Sprintf("(%s) Error threshold met", ctx.Value("profile-name"))
			if totalRequests >= 10 {
				switch separated[1] {
				case "raw":
					if failedRequests > threshold {
						common.GetLogger().Log(errMsg)
						cancelCtx()
					}
				default:
					errPercent := failedRequests * 100 / totalRequests

					if errPercent > threshold {
						common.GetLogger().Log(errMsg)
						cancelCtx()
					}
				}
			}

		}
	}
}
