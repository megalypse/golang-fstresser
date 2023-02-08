package customprofile

import (
	"context"
	"sync"
)

var defaultReqWg sync.WaitGroup

func deployDefaultRequester(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	csp *CustomStressProfile,
	loadsConsumer <-chan DefaultRequesterPayload,
	reqCountProducer chan<- int,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case load := <-loadsConsumer:
			for i := 0; i < load.Rps; i++ {
				go func() {
					res := csp.MakeRequestUsecase.Request(cancelCtx, load.Request, csp.Config.GlobalHeaders)
					reqCountProducer <- res.StatusCode
				}()
			}
		}
	}
}

func deployCustomRequester(
	ctx context.Context,
	cancelCtx context.CancelFunc,
	csp *CustomStressProfile,
	loadsConsumer <-chan CustomRequesterPayload,
	reqCountProducer chan<- int,
) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			defaultReqWg.Wait()
			return
		case load := <-loadsConsumer:
			for i := 0; i < load.CustomLoadConfig.Rps; i++ {
				defaultReqWg.Add(1)
				go func() {
					res := csp.MakeRequestUsecase.Request(cancelCtx, load.Request, csp.Config.GlobalHeaders)
					reqCountProducer <- res.StatusCode
					defaultReqWg.Done()
				}()
			}
		}
	}
}
