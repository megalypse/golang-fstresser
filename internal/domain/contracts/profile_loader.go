package contracts

type ProfileLoader interface {
	LoadProfile() []StressProfile
}
