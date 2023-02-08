package main

import (
	"context"
	"log"
	"os"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/main/factory"
)

func main() {
	// runtime.GOMAXPROCS(6)
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	defer common.HandlePanic(ctx, cancelCtx)

	path := os.Getenv("FSTRESSER_PROFILES_PATH")
	if path == "" {
		log.Fatal("Profiles path not defined")
	}

	loader := factory.MakeLocalProfileLoader()

	profiles := loader.LoadProfile(ctx, cancelCtx, path)

	for _, profile := range profiles {
		profile.StartLoad(ctx, cancelCtx)
	}
}
