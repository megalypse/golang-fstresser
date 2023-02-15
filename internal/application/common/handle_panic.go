package common

import (
	"context"
	"fmt"
)

func HandlePanic(ctx context.Context, cancelCtx context.CancelFunc) {
	if err := recover(); err != nil {
		errMsg := fmt.Sprintf("Runtime panicked with %v", err)

		GracefulVarnish(ctx, cancelCtx, errMsg)

		GetLogger(ctx).RegisterLogs()
	}
}
