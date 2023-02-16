package common

import "context"

// A convenient funcion to do a final logging before canceling a context
func LogFinale(ctx context.Context, cancelCtx context.CancelFunc, message string) {
	GetLogger(ctx).Log(message)
	cancelCtx()
}
