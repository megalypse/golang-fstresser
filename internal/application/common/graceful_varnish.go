package common

import "context"

func GracefulVarnish(cancelCtx context.CancelFunc, message string) {
	GetLogger().Log(message)
	cancelCtx()
}
