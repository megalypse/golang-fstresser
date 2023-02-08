package common

import "context"

func GracefulVarnish(ctx context.Context, cancelCtx context.CancelFunc, message string) {
	GetLogger(ctx).Log(message)
	cancelCtx()
}
