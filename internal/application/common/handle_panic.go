package common

import (
	"context"
	"fmt"
)

// HandlePanic logs the panic message, cancel the context through `LogFinale`
// and finally register the logs before proceeding to shuting down the execution.
func HandlePanic(ctx context.Context, cancelCtx context.CancelFunc) {
	if err := recover(); err != nil {
		errMsg := fmt.Sprintf("Runtime panicked with %v", err)

		LogFinale(ctx, cancelCtx, errMsg)

		GetLogger(ctx).RegisterLogs()
	}
}
