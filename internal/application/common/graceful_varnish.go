package common

import "context"

// TODO: change this function name

// A convenient funcion to do a final logging before canceling a context
func GracefulVarnish(ctx context.Context, cancelCtx context.CancelFunc, message string) {
	GetLogger(ctx).Log(message)
	cancelCtx()
}
