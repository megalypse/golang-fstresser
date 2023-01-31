package usecase

import "context"

type ProfileLoader interface {
	LoadProfile(cancelCtx context.CancelFunc, profilesPath string) []StressProfile
}
