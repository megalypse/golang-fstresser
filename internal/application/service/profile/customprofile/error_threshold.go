package customprofile

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/megalypse/golang-fstresser/internal/application/common"
)

func deployErrorThresholdAnalyzer(
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
	treshold, err := strconv.Atoi(separated[0])
	if err != nil {
		common.GetLogger().Log(err.Error())
		cancelCtx()
		return
	}

	for {
		select {
		case <-ctx.Done():
			message := fmt.Sprintf("Requests done: %d\nSuccess: %d\nFailed: %d", totalRequests, successfullRequests, failedRequests)
			common.GetLogger().Log(message)
			return
		case httpStatus := <-reqCountConsumer:
			if httpStatus >= 200 && httpStatus < 300 {
				successfullRequests++
			} else {
				failedRequests++
			}

			totalRequests = successfullRequests + failedRequests

			switch separated[1] {
			case "raw":
				if failedRequests > treshold {
					common.GetLogger().Log("Error threshold met.")
					cancelCtx()
				}
			default:
				errPercent := failedRequests * 100 / totalRequests

				if errPercent > treshold {
					common.GetLogger().Log("Error threshold met.")
					cancelCtx()
				}
			}
		}
	}
}
