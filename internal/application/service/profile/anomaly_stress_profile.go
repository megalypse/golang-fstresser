package profile

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/service"
	"github.com/megalypse/golang-fstresser/internal/application/service/logger"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

var mu sync.Mutex
var wg sync.WaitGroup
var lgr logger.Logger

func init() {
	lgr = logger.NewLogger()
}

type AnomalyStressProfile struct {
	RequestService entity.RequestService
	Req            entity.Request
	Config         Config
	State          State
}

func (asp *AnomalyStressProfile) bootProfile() {
	asp.Config.ExpectedExecutionTime = asp.Config.BeginAnomalyAfter + asp.Config.AnomalyDuration + asp.Config.HoldPeakAfterAnomalyFor
	asp.Config.ExpectedAnomalyDeadline = asp.Config.BeginAnomalyAfter + asp.Config.AnomalyDuration
	asp.Config.RampUpPace = float64(asp.Config.PeakRps) / asp.Config.RampUpTime.Seconds()

	asp.State.CurrentRps = func() float64 {
		if asp.Config.RampUpPace > 1 {
			return asp.Config.RampUpPace
		}

		return 1
	}()

	asp.State.EffectiveRps = int64(asp.State.CurrentRps)
	asp.State.IsDefaultFluxActive = true
	asp.State.IsAnomalyFluxActive = false

	message := (fmt.Sprintf(
		"\n==============================\nExpected execution time: %v\nExpected anomaly deadline: %v\nRampup pace: %f\nInitial Rps: %d\n==============================\n",
		asp.Config.ExpectedExecutionTime,
		asp.Config.ExpectedAnomalyDeadline,
		asp.Config.RampUpPace,
		asp.State.EffectiveRps,
	))

	lgr.Log(message)
}

type Config struct {
	PeakRps int64

	// Amount of time to reach peak RPS (relative to minute zero)
	RampUpTime time.Duration

	// Amount of time to wait before anomaly starts (relative to minute zero)
	BeginAnomalyAfter time.Duration

	// For how long the anomaly will be sustained
	AnomalyDuration time.Duration

	// The anomaly Rps will be calculated by (peakRps * AnomalyRps)
	AnomalyRps int

	// For how long the peak rps will be held after the anomaly ends
	HoldPeakAfterAnomalyFor time.Duration

	// Expected execution time calculated by the following formulae:
	// BeginAnomalyAfter + AnomalyDuration + PeakAfterAnomalyHold
	ExpectedExecutionTime time.Duration

	// Computed value for when the anomaly should end. (BeginAnomalyAfter + AnomalyDuration)
	ExpectedAnomalyDeadline time.Duration

	// Computed value for how many Rps current pace will get an increase of per minute
	RampUpPace float64
}

type State struct {
	CurrentRps          float64
	EffectiveRps        int64
	ComparableRps       int
	Runtime             time.Duration
	IsDefaultFluxActive bool
	IsAnomalyFluxActive bool
}

func (asp *AnomalyStressProfile) StartLoad() {
	asp.bootProfile()

	ctx := context.Background()
	ctx, cancelContext := context.WithCancel(ctx)
	defer cancelContext()

	wg.Add(1)
	go deployOrchestrator(ctx, cancelContext, time.Second, asp)
	go deployDefaultFlux(ctx, asp)
	go deployAnomalyFlux(ctx, asp)

	wg.Wait()
}

func deployOrchestrator(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	durationMetric time.Duration,
	asp *AnomalyStressProfile,
) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()

			logsDir := "../../../logs"
			os.Mkdir(logsDir, 0777)

			fileName := fmt.Sprintf(logsDir+"/%d.txt", time.Now().UnixMilli())
			err := os.WriteFile(fileName, lgr.GetBuffer(), 0644)

			if err != nil {
				log.Fatal(err)
			}

			return
		default:
			time.Sleep(durationMetric)
			mu.Lock()

			asp.State.Runtime += durationMetric
			currentRuntime := asp.State.Runtime

			if currentRuntime == asp.Config.ExpectedExecutionTime {
				cancelCtx()
			}

			isGreaterThanAnomalyInterval := currentRuntime >= asp.Config.BeginAnomalyAfter
			isLowerThanAnomalyInterval := currentRuntime <= asp.Config.ExpectedAnomalyDeadline
			isAnomalyInterval := isGreaterThanAnomalyInterval && isLowerThanAnomalyInterval

			prevEffectiveRps := asp.State.EffectiveRps
			if !isAnomalyInterval && asp.State.EffectiveRps < asp.Config.PeakRps {
				asp.State.CurrentRps = asp.State.CurrentRps + asp.Config.RampUpPace
				asp.State.EffectiveRps = int64(asp.State.CurrentRps)
				asp.State.ComparableRps = int(asp.State.EffectiveRps)
			}

			if prevEffectiveRps != asp.State.EffectiveRps {
				if isAnomalyInterval {
					message := fmt.Sprintf("\nRuntime: %f\nRPS: %d (ANOMALY)\n", currentRuntime.Seconds(), asp.Config.AnomalyRps)
					lgr.Log(message)
				} else {
					message := fmt.Sprintf("\nRuntime: %f\nRPS: %d\n", currentRuntime.Seconds(), asp.State.EffectiveRps)
					lgr.Log(message)
				}
			}

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
		mu.Lock()
		isDefaultFluxActive := asp.State.IsDefaultFluxActive
		mu.Unlock()

		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			if isDefaultFluxActive {
				mu.Lock()
				currentRps := asp.State.ComparableRps
				mu.Unlock()

				for i := 0; i < currentRps; i++ {
					go service.MakeRequest(&asp.Req, asp.RequestService)
				}
			}
			time.Sleep(time.Second)
		}
	}
}

func deployAnomalyFlux(ctx context.Context, asp *AnomalyStressProfile) {
	wg.Add(1)

	for {
		mu.Lock()
		isAnomalyFluxActive := asp.State.IsAnomalyFluxActive
		mu.Unlock()

		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			if isAnomalyFluxActive {
				mu.Lock()
				anomalyRps := asp.Config.AnomalyRps
				mu.Unlock()

				for i := 0; i < anomalyRps; i++ {
					go service.MakeRequest(&asp.Req, asp.RequestService)
				}
			}
			time.Sleep(time.Second)
		}
	}
}
