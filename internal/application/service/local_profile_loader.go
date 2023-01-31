package service

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/application/service/profile/customprofile"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
)

type profilesWrapper struct {
	Profiles []customprofile.CustomStressProfile
}

type LocalProfileLoader struct {
	MakeRequestUsecase usecase.MakeRequestUsecase
}

func (lpl LocalProfileLoader) LoadProfile(cancelCtx context.CancelFunc, profilesPath string) []usecase.StressProfile {
	if profilesPath == "" {
		common.GracefulVarnish(cancelCtx, "Profiles path not provided. Ending execution...")
	}

	result, err := os.ReadFile(profilesPath)

	if err != nil {
		common.GracefulVarnish(cancelCtx, err.Error())
	}

	return lpl.objectifyProfiles(result)
}

func (lpl LocalProfileLoader) objectifyProfiles(bytes []byte) []usecase.StressProfile {
	holder := new(profilesWrapper)

	err := json.Unmarshal(bytes, holder)
	if err != nil {
		log.Fatal(err.Error())
	}

	profiles := make([]usecase.StressProfile, 0, len(holder.Profiles))
	for _, v := range holder.Profiles {
		v.MakeRequestUsecase = lpl.MakeRequestUsecase

		profiles = append(profiles, v)
	}

	return profiles
}
