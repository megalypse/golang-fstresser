package service

import (
	"context"
	"encoding/json"
	"os"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/contracts"
)

type ProfilesWrapper struct {
	Profiles []contracts.AnomalyStressProfile
}

type LocalProfileLoader struct{}

func (LocalProfileLoader) LoadProfile(cancelCtx context.CancelFunc) []contracts.AnomalyStressProfile {
	profilesPath := os.Getenv("FSTRESSER_PROFILES_PATH")

	if profilesPath == "" {
		common.GracefulVarnish(cancelCtx, "Profiles path not provided. Finishing execution...")
	}

	result, err := os.ReadFile(profilesPath)
	if err != nil {
		common.GracefulVarnish(cancelCtx, err.Error())
	}

	return parseJsonProfile(result)
}

func parseJsonProfile(bytes []byte) []contracts.AnomalyStressProfile {
	holder := new(ProfilesWrapper)

	json.Unmarshal(bytes, holder)

	return holder.Profiles
}
