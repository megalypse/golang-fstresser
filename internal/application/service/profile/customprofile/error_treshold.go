package customprofile

import (
	"context"
	"strconv"
	"strings"

	"github.com/megalypse/golang-fstresser/internal/application/common"
)

func errorTreshold(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	httpStatusConsumer <-chan int,
	cpc *CustomProfileConfig,
) {
	var successfullRequests int
	var failedRequests int
	var totalRequests int

	for {
		select {
		case <-ctx.Done():
			return
		case httpStatus := <-httpStatusConsumer:
			if httpStatus >= 200 && httpStatus < 300 {
				successfullRequests++
			} else {
				failedRequests++
			}

			totalRequests = successfullRequests + failedRequests

			separated := strings.Split(cpc.ErrorTreshold, ":")
			treshold, err := strconv.Atoi(separated[0])
			if err != nil {
				common.GetLogger().Log(err.Error())
				cancelCtx()
				return
			}

			switch separated[1] {
			case "raw":
				if failedRequests > treshold {
					common.GetLogger().Log("Error treshold met.")
					cancelCtx()
				}
			default:
				errPercent := failedRequests * 100 / totalRequests

				if errPercent > treshold {
					common.GetLogger().Log("Error treshold met.")
					cancelCtx()
				}
			}
		}
	}
}
