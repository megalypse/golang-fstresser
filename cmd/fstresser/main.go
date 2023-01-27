package main

import (
	"context"

	"github.com/megalypse/golang-fstresser/internal/application/service"
)

func main() {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	loader := service.LocalProfileLoader{}

	profile := loader.LoadProfile(cancelCtx, "/Users/megalypse/Documents/Projects/go/fstresser/resources/first_profile.json")[0]

	profile.StartLoad(ctx, cancelCtx)
}
