package main

import (
	"context"
	"runtime"

	"github.com/megalypse/golang-fstresser/internal/application/service"
)

func main() {
	runtime.GOMAXPROCS(6)
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	loader := service.LocalProfileLoader{}

	profiles := loader.LoadProfile(cancelCtx, "/Users/megalypse/Documents/Projects/go/fstresser/resources/first_profile.json")

	for _, profile := range profiles {
		if profile.IsActive {
			profile.StartLoad(ctx, cancelCtx)
		}
	}
}
