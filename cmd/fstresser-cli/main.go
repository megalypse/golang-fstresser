package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/megalypse/golang-fstresser/internal/application/service"
	"github.com/megalypse/golang-fstresser/internal/application/service/profile/customprofile"
	"github.com/megalypse/golang-fstresser/internal/main/factory"
)

var wg sync.WaitGroup

func main() {
	path := os.Getenv("FSTRESSER_PROFILES_PATH")

	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	if path == "" {
		log.Fatal("Profiles path not defined")
	}

	loader := factory.MakeLocalProfileLoader()

	profiles := findProfiles(path)
	wrappers := make([]*service.ProfilesWrapper, 0, len(profiles))

	for _, v := range profiles {
		bytes, err := os.ReadFile(v)
		if err != nil {
			log.Fatal(err.Error())
		}

		wrapper := loader.(service.LocalProfileLoader).ObjectifyProfiles(bytes)
		wrappers = append(wrappers, wrapper)
	}

	fmt.Println("Choose the desired profiles to be run:")
	for i, v := range wrappers {
		fmt.Println(i, v.ProfileName)
	}

	var chosenProfilesRaw string
	fmt.Scan(&chosenProfilesRaw)

	splittedChoices := strings.Split(chosenProfilesRaw, ",")

	for _, v := range splittedChoices {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, profile := range wrappers[parsed].Profiles {
			profile.MakeRequestUsecase = loader.(service.LocalProfileLoader).MakeRequestUsecase

			newCtx := context.WithValue(ctx, "profile-name", wrappers[parsed].ProfileName)
			newCtx, cancelNewCtx := context.WithCancel(newCtx)

			wg.Add(1)
			go func(profile customprofile.CustomStressProfile) {
				profile.StartLoad(newCtx, cancelNewCtx)
				wg.Done()
			}(profile)
		}
	}

	wg.Wait()
}

func findProfiles(path string) []string {
	profilesPath := os.Getenv("FSTRESSER_PROFILES_PATH")
	entries, err := os.ReadDir(profilesPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	profiles := make([]string, 0, len(entries))
	for _, v := range entries {
		profiles = append(profiles, path+v.Name())
	}

	return profiles
}
