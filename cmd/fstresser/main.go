package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/megalypse/golang-fstresser/internal/application/service"
)

func main() {
	runtime.GOMAXPROCS(6)
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	loader := service.LocalProfileLoader{}

	path := os.Getenv("FSTRESSER_PROFILES_PATH")
	if path == "" {
		log.Fatal("Profiles path not defined")
	}

	profiles := loader.LoadProfile(cancelCtx, path)

	for _, profile := range profiles {
		if profile.IsActive {
			profile.StartLoad(ctx, cancelCtx)
		}
	}
}
