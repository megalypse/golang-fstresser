package service

import (
	"context"
	"sync"
	"time"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

var mu sync.Mutex
var wg sync.WaitGroup
var httpService HttpService

func init() {
	httpService = HttpService{}
}

type AnomalyStressProfile struct {
	Req    entity.Request
	Config *config
	State  state
}

type config struct {
	PeakRps int

	// Amount of time to reach peak RPS (relative to minute zero)
	RampUpTime int

	// Amount of time to wait before anomaly starts (relative to minute zero)
	BeginAnomalyAfter int

	// For how long the anomaly will be sustained
	AnomalyDuration int

	// The anomaly Rps will be calculated by (peakRps * AnomalyMultiplier)
	AnomalyMultiplier int

	// For how long the peak rps will be held after the anomaly ends
	PeakAfterAnomalyHold int

	// Expected execution time calculated by the following formulae:
	// BeginAnomalyAfter + AnomalyDuration + PeakAfterAnomalyHold
	ExpectedExecutionTime int

	// Computed value for when the anomaly should end. (BeginAnomalyAfter + AnomalyDuration)
	ExpectedAnomalyDeadline int

	// Computed value for how many Rps current pace will get an increase of per minute
	RampUpPace int
}

func (asp *AnomalyStressProfile) setComputedValues() {
	asp.Config.ExpectedExecutionTime = asp.Config.RampUpTime + asp.Config.AnomalyDuration + asp.Config.PeakAfterAnomalyHold
	asp.Config.ExpectedAnomalyDeadline = asp.Config.BeginAnomalyAfter + asp.Config.AnomalyDuration
	asp.Config.RampUpPace = asp.Config.PeakRps / asp.Config.RampUpTime
	asp.State.CurrentRps = asp.Config.RampUpPace
}

type state struct {
	CurrentRps          int
	Runtime             int
	IsDefaultFluxActive bool
	IsAnomalyFluxActive bool
}

func (asp *AnomalyStressProfile) StartLoad() {
	asp.setComputedValues()
	expectedExecutionTime := asp.Config.ExpectedExecutionTime

	ctx := context.Background()
	ctx, cancelContext := context.WithTimeout(ctx, time.Minute*time.Duration(expectedExecutionTime))
	defer cancelContext()

	ticker := time.NewTicker(time.Minute)

	wg.Add(1)
	go chronos(ctx, ticker, asp)
	go deployDefaultFlux(ctx, asp)
	go deployAnomalyFlux(ctx, asp)

	wg.Wait()
}

func chronos(
	ctx context.Context,
	ticker *time.Ticker,
	asp *AnomalyStressProfile,
) {
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			wg.Done()
			return
		case <-ticker.C:
			mu.Lock()

			asp.State.Runtime++
			currentRuntime := asp.State.Runtime

			isGreaterThanAnomalyInterval := currentRuntime >= asp.Config.BeginAnomalyAfter
			isLowerThanAnomalyInterval := currentRuntime <= asp.Config.ExpectedAnomalyDeadline
			isAnomalyInterval := isGreaterThanAnomalyInterval && isLowerThanAnomalyInterval

			if asp.State.IsDefaultFluxActive == asp.State.IsAnomalyFluxActive {
				asp.State.IsDefaultFluxActive = true
				asp.State.IsAnomalyFluxActive = false
			}

			if asp.State.IsDefaultFluxActive && !asp.State.IsAnomalyFluxActive && isAnomalyInterval {
				asp.State.IsDefaultFluxActive = false
				asp.State.IsAnomalyFluxActive = true
			}

			if !asp.State.IsDefaultFluxActive && asp.State.IsAnomalyFluxActive && !isAnomalyInterval {
				asp.State.IsDefaultFluxActive = true
				asp.State.IsAnomalyFluxActive = false
			}

			mu.Unlock()
		}
	}
}

func deployDefaultFlux(ctx context.Context, asp *AnomalyStressProfile) {
	wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			mu.Lock()
			isDefaultFluxActive := asp.State.IsDefaultFluxActive
			currentRps := asp.State.CurrentRps
			mu.Unlock()

			if isDefaultFluxActive {
				for i := 0; i < currentRps; i++ {
					go MakeRequest(&asp.Req, httpService)
				}

				time.Sleep(time.Second)
			}
		}
	}
}

func deployAnomalyFlux(ctx context.Context, asp *AnomalyStressProfile) {
	wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			mu.Lock()
			isAnomalyFluxActive := asp.State.IsAnomalyFluxActive
			anomalyRps := asp.State.CurrentRps * asp.Config.AnomalyMultiplier
			mu.Unlock()

			if isAnomalyFluxActive {
				for i := 0; i < anomalyRps; i++ {
					go MakeRequest(&asp.Req, httpService)
				}

				time.Sleep(time.Second)
			}
		}
	}
}
