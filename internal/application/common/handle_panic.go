package common

import (
	"context"
	"fmt"
)

func HandlePanic(cancelCtx context.CancelFunc) {
	if err := recover(); err != nil {
		errMsg := fmt.Sprint(err)
		GetLogger().Log(errMsg)

		cancelCtx()

		GetLogger().RegisterLogs()
	}
}
