package main

import (
	"context"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/main/factory"
)

func main() {
	common.SetMaxProcs()
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)

	defer common.HandlePanic(ctx, cancelCtx)

	loader := factory.MakeLocalProfileLoader()

	profiles := loader.LoadProfile(ctx, cancelCtx)

	for _, profile := range profiles {
		profile.StartLoad(ctx, cancelCtx)
	}
}
