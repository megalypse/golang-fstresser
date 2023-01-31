package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/megalypse/golang-fstresser/internal/main/factory"
)

func main() {
	runtime.GOMAXPROCS(6)
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	path := os.Getenv("FSTRESSER_PROFILES_PATH")
	if path == "" {
		log.Fatal("Profiles path not defined")
	}

	loader := factory.MakeLocalProfileLoader()

	profiles := loader.LoadProfile(cancelCtx, path)

	for _, profile := range profiles {
		profile.StartLoad(ctx, cancelCtx)
	}
}
