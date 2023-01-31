package usecase

import "context"

type StressProfile interface {
	StartLoad(context.Context, context.CancelFunc)
}
