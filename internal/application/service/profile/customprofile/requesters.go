package customprofile

import (
	"context"

	"github.com/megalypse/golang-fstresser/internal/application/common"
)

func deployDefaultRequester(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	loadsConsumer <-chan DefaultRequesterPayload,
	reqCountProducer chan<- int,
) {
	wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case load := <-loadsConsumer:
			for i := 0; i < load.Rps; i++ {
				go func() {
					res := common.MakeLightweightRequest(cancelCtx, load.Request)
					reqCountProducer <- res.StatusCode
				}()
			}
		}
	}
}

func deployCustomRequester(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	loadsConsumer <-chan CustomRequesterPayload,
	reqCountProducer chan<- int,
) {
	wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case load := <-loadsConsumer:
			for i := 0; i < load.CustomLoadConfig.Rps; i++ {
				go func() {
					res := common.MakeLightweightRequest(cancelCtx, load.Request)
					reqCountProducer <- res.StatusCode
				}()
			}
		}
	}
}
