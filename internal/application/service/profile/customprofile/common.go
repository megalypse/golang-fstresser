package customprofile

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

var wg sync.WaitGroup
var requestQueue []*entity.Request

func init() {
	requestQueue = make([]*entity.Request, 0)
}

// logRps contains common logic that got extracted to avoid repetition.
//
// It only logs the provided message if `prevRps` != `currentRps`. This means the message will only be logged when
// the current RPS changes.
func logRps(ctx context.Context, prevRps, currentRps int, runtime time.Duration) {
	if prevRps != currentRps {
		common.GetLogger(ctx).Log(fmt.Sprintf("(%s) Runtime: %s, Rps: %d", ctx.Value(common.GetCtxKey("profile-name")), runtime.String(), currentRps))
	}
}

// isCustomLoadWindow receives all the custom loads from the profile, and check
// the custom load time window against `now`. If `now` is in the interval, it means
// the current timestamp should be a custom load.
func isCustomLoadWindow(cpc *CustomProfileConfig, now int64) *CustomLoad {
	for _, v := range cpc.CustomLoads {
		if now >= v.GetStartPoint() && now < v.GetEndPoint() {
			return &v
		}
	}

	return nil
}

/*
makeRequestQueueCounter creates an instance of the iterator to be used to keep track of
which request should be sent.
*/
func makeRequestQueueCounter(maxCount int) RateCounter {
	if maxCount < 0 {
		log.Fatal("Max count must be a positive number")
	}

	return RateCounter{
		maxValue:     maxCount,
		currentValue: -1,
	}
}

/*
prepareRequests convert the structs from the profile to the request entity to be used with
MakeLightweightRequest.
*/
func prepareRequests(csp *CustomStressProfile) {
	for _, v := range csp.Requests {
		requestEntity := v.ToEntity()

		for i := 0; i < v.Rate; i++ {
			requestQueue = append(requestQueue, &requestEntity)
		}
	}
}

// calculateCustomLoadsWindows receives the profile and calculate the timewindow
// for all the custom loads inside it.
func calculateCustomLoadsWindows(startTime time.Time, csp *CustomStressProfile) {
	for i, v := range csp.Config.CustomLoads {
		v.calculateWindow(startTime)

		csp.Config.CustomLoads[i] = v
	}
}

// getInitialRps makes sure the initial RPS will be at least 1
func getInitialRps(cpc *CustomProfileConfig) float64 {
	tempRps := cpc.RpsRampupPace

	if tempRps > 1 {
		return tempRps
	}

	return 1
}
