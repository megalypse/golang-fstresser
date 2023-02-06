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

type ProfilesWrapper struct {
	ProfileName string
	Profiles    []customprofile.CustomStressProfile
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

	profilesWrapper := lpl.ObjectifyProfiles(result)

	profiles := make([]usecase.StressProfile, 0, len(profilesWrapper.Profiles))
	for _, v := range profilesWrapper.Profiles {
		v.MakeRequestUsecase = lpl.MakeRequestUsecase

		profiles = append(profiles, v)
	}

	return profiles
}

func (lpl LocalProfileLoader) ObjectifyProfiles(bytes []byte) *ProfilesWrapper {
	holder := new(ProfilesWrapper)

	err := json.Unmarshal(bytes, holder)
	if err != nil {
		log.Fatal(err.Error())
	}

	return holder
}
