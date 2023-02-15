package customprofile

import (
	"context"
	"time"

	"github.com/megalypse/golang-fstresser/internal/application/common"
)

// This one purpose is to keep track of whenever a rampup should be made
func deployRpsComposer(ctx context.Context, cancelCtx context.CancelFunc, startTime time.Time, cpc *CustomProfileConfig) chan int {
	rpsChan := make(chan int)
	ticker := time.NewTicker(cpc.RpsIncreaseInterval.Duration)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer common.HandlePanic(ctx, cancelCtx)

		rawRps := getInitialRps(cpc)
		effectiveRps := int(rawRps)
		tickerDeadline := startTime.Add(cpc.RampUpTime.Duration).Unix()

	l1:
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// In case the desired peak RPS have not been met at the end of the rampup time
				nowUnix := time.Now().Unix()
				if nowUnix > tickerDeadline {
					if effectiveRps != cpc.PeakRps {
						rpsChan <- cpc.PeakRps
					}

					common.GetLogger(ctx).Log("Rampup finished")
					break l1
				} else {
					rawRps += cpc.RpsRampupPace
					effectiveRps = int(rawRps)

					if effectiveRps != 0 {
						rpsChan <- effectiveRps
					}
				}
			}

		}
	}()

	return rpsChan
}
