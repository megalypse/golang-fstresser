package customprofile

import (
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

func logRps(prevRps, currentRps int, runtime int64) {
	if prevRps != currentRps {
		common.GetLogger().Log(fmt.Sprintf("Runtime: %ds, Rps: %d", runtime, currentRps))
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