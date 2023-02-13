package customprofile

import (
	"context"
)

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
		case load, ok := <-loadsConsumer:
			if !ok {
				return
			}

			for i := 0; i < load.Rps; i++ {
				go func(isChannelOpen bool) {
					res := csp.MakeRequestUsecase.Request(ctx, cancelCtx, load.Request, csp.Config.GlobalHeaders)

					if isChannelOpen {
						reqCountProducer <- res.StatusCode
					}
				}(ok)
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
			return
		case load, ok := <-loadsConsumer:
			if !ok {
				return
			}

			for i := 0; i < load.CustomLoadConfig.Rps; i++ {
				go func(isChannelOpen bool) {
					res := csp.MakeRequestUsecase.Request(ctx, cancelCtx, load.Request, csp.Config.GlobalHeaders)

					if isChannelOpen {
						reqCountProducer <- res.StatusCode
					}
				}(ok)
			}
		}
	}
}
