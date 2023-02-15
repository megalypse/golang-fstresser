package usecase

import (
	"context"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type MakeRequestUsecase interface {
	Request(context.Context, context.CancelFunc, *entity.Request, map[string]string) *entity.Response
}