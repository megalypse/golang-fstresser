package service

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/application/service/profile/customprofile"
)

type profilesWrapper struct {
	Profiles []customprofile.CustomStressProfile
}

type LocalProfileLoader struct{}

func (LocalProfileLoader) LoadProfile(cancelCtx context.CancelFunc, profilesPath string) []customprofile.CustomStressProfile {
	if profilesPath == "" {
		common.GracefulVarnish(cancelCtx, "Profiles path not provided. Ending execution...")
	}

	result, err := os.ReadFile(profilesPath)

	if err != nil {
		common.GracefulVarnish(cancelCtx, err.Error())
	}

	return ObjectifyProfiles(result)
}

func ObjectifyProfiles(bytes []byte) []customprofile.CustomStressProfile {
	holder := new(profilesWrapper)

	err := json.Unmarshal(bytes, holder)
	if err != nil {
		log.Fatal(err.Error())
	}

	return holder.Profiles
}
