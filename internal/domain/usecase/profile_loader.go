package usecase

import "context"

type ProfileLoader interface {
	LoadProfile(ctx context.Context, cancelCtx context.CancelFunc) []StressProfile
}
