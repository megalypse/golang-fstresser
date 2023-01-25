package contracts

type ProfileLoader interface {
	LoadProfile() []AnomalyStressProfile
}
