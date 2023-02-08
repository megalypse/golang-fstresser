package common

import (
	"context"
	"fmt"
)

func HandlePanic(ctx context.Context, cancelCtx context.CancelFunc) {
	if err := recover(); err != nil {
		errMsg := fmt.Sprint(err)
		GetLogger(ctx).Log(errMsg)

		cancelCtx()

		GetLogger(ctx).RegisterLogs()
	}
}
