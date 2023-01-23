package profile

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

var wg sync.WaitGroup
var lgr *common.Logger

func init() {
	lgr = common.GetLogger()
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
	CurrentRps    float64
	EffectiveRps  int64
	ComparableRps int
	Runtime       time.Duration
}

func (asp *AnomalyStressProfile) StartLoad() {
	asp.bootProfile()

	ctx := context.Background()
	ctx, cancelContext := context.WithCancel(ctx)
	defer cancelContext()

	deployOrchestrator(ctx, cancelContext, asp)
}

func deployOrchestrator(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	asp *AnomalyStressProfile,
) {
	startPoint := time.Now()
	anomalyStartsAt := startPoint.Add(asp.Config.BeginAnomalyAfter)
	anomalyEndsAt := anomalyStartsAt.Add(asp.Config.AnomalyDuration)
	expectedDeadline := anomalyEndsAt.Add(asp.Config.HoldPeakAfterAnomalyFor)

	defaultFluxChan := make(chan int)
	anomalyFluxChan := make(chan int)
	defer close(defaultFluxChan)
	defer close(anomalyFluxChan)

	currentRps := 0.0
	effectiveRps := 1
	previousEffectiveRps := 1

	go deployDefaultFlux(ctx, cancelCtx, &asp.Req, defaultFluxChan)
	go deployAnomalyFlux(ctx, cancelCtx, &asp.Req, anomalyFluxChan)

l1:
	for {

		select {
		case <-ctx.Done():
			wg.Done()

			lgr.RegisterLogs()

			break l1
		default:
			currentRps += asp.Config.RampUpPace
			previousEffectiveRps = effectiveRps
			effectiveRps = int(currentRps)

			now := time.Now()
			runtime := now.Unix() - startPoint.Unix()

			if now.Unix() > expectedDeadline.Unix() {
				cancelCtx()
				continue
			}

			if now.Unix() < anomalyStartsAt.Unix() || now.Unix() > anomalyEndsAt.Unix() {
				if effectiveRps != previousEffectiveRps {
					lgr.Log(fmt.Sprintf("\nRuntime: %ds\nRPS: %d\n", runtime, effectiveRps))
				}

				defaultFluxChan <- effectiveRps
			} else {
				lgr.Log(fmt.Sprintf("\nRuntime: %ds\nRPS: %d (ANOMALY)\n", runtime, effectiveRps))
				anomalyFluxChan <- asp.Config.AnomalyRps
			}

			time.Sleep(time.Second)
		}
	}

	wg.Wait()
}

func deployDefaultFlux(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	req *entity.Request,
	fluxChan <-chan int,
) {
	wg.Add(1)
	defer wg.Done()

l1:
	for {
		select {
		case <-ctx.Done():
			break l1
		case rps := <-fluxChan:
			go func() {
				for i := 0; i < rps; i++ {
					go common.MakeLightweightRequest(cancelCtx, req)
				}
			}()
		}
	}
}

func deployAnomalyFlux(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	req *entity.Request,
	fluxChan <-chan int,
) {
	wg.Add(1)
	defer wg.Done()

l1:
	for {
		select {
		case <-ctx.Done():
			break l1
		case rps := <-fluxChan:
			go func() {
				for i := 0; i < rps; i++ {
					go common.MakeLightweightRequest(cancelCtx, req)
				}
			}()
		}
	}
}
