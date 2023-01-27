package customprofile

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testWg sync.WaitGroup

func TestRpsComposer(t *testing.T) {
	assert := assert.New(t)

	testWg.Add(1)
	go func() {
		ctx := context.Background()
		ctx, cancelCtx := context.WithCancel(ctx)
		rpsChan, rpsCtx := makeRpsChan(ctx)
		currentRps := 0

	l1:
		for {
			select {
			case <-rpsCtx.Done():
				assert.Equal(10, currentRps)
				break l1
			case rps := <-rpsChan:
				if rps != 0 {
					currentRps = rps
				}
			}
		}

		cancelCtx()
		testWg.Done()
	}()

	testWg.Add(1)
	go func() {
		ctx := context.Background()
		ctx, cancelCtx := context.WithCancel(ctx)
		rpsChan, rpsCtx := makeRpsChan(ctx)
		currentRps := 0

		go func() {
			time.Sleep(time.Millisecond * 5100)

			cancelCtx()
		}()

	l1:
		for {
			select {
			case <-rpsCtx.Done():
				assert.Equal(6, currentRps)
				break l1
			case rps := <-rpsChan:
				if rps != 0 {
					currentRps = rps
				}
			}
		}

		cancelCtx()
		testWg.Done()
	}()

	testWg.Wait()
}

func makeRpsChan(ctx context.Context) (<-chan int, context.Context) {
	now := time.Now()
	config := CustomProfileConfig{
		RampUpTime:          DurationInput{time.Second * 10},
		RpsIncreaseInterval: DurationInput{time.Second},
		PeakRps:             10,
	}

	config.bootstrap()

	return deployRpsComposer(ctx, now, &config)
}
