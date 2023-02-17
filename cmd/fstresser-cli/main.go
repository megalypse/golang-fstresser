package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/application/service"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
	"github.com/megalypse/golang-fstresser/internal/main/factory"
)

var wg sync.WaitGroup

func main() {
	common.SetMaxProcs()

	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	defer cancelCtx()
	defer common.HandlePanic(ctx, cancelCtx)

	path := common.GetResourcesPath()

	loader := factory.MakeLocalProfileLoader()
	wrappers := getWrappers(loader, path)
	indexes := chooseProfilesIndexes(wrappers)

	runProfiles(ctx, loader, indexes, wrappers)

	wg.Wait()
}

// runProfiles run the user selected profiles.
// Each profile will run on its own go routine.
func runProfiles(
	ctx context.Context,
	loader usecase.ProfileLoader,
	indexes []int,
	wrappers []*service.ProfilesWrapper,
) {
	for _, v := range indexes {
		for _, profile := range wrappers[v].Profiles {
			profile.MakeRequestUsecase = loader.(service.LocalProfileLoader).MakeRequestUsecase

			ctx = context.WithValue(ctx, common.GetCtxKey("profile-name"), wrappers[v].ProfileName)
			ctx, cancelNewCtx := context.WithCancel(ctx)

			wg.Add(1)
			go func(profile usecase.StressProfile) {
				profile.StartLoad(ctx, cancelNewCtx)
				wg.Done()
			}(profile)
		}
	}
}

/*
getWrappers uses the given path and the loader to read all the profiles there located
*/
func getWrappers(loader usecase.ProfileLoader, path string) []*service.ProfilesWrapper {
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

	return wrappers
}

// An aux function to extract the cli logic from the main workflow
func chooseProfilesIndexes(wrappers []*service.ProfilesWrapper) []int {
	fmt.Println("Choose the desired profiles to be run:")
	for i, v := range wrappers {
		fmt.Println(i, v.ProfileName)
	}

	var chosenProfilesRaw string
	fmt.Scan(&chosenProfilesRaw)

	splittedChoices := strings.Split(chosenProfilesRaw, ",")

	finalList := make([]int, 0, len(splittedChoices))

	for _, v := range splittedChoices {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err.Error())
		}

		finalList = append(finalList, parsed)
	}

	return finalList
}

// Find all the stress profile json files in the given path
func findProfiles(path string) []string {
	profilesPath := common.GetResourcesPath()
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
