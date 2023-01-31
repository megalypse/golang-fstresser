package usecase

import (
	"context"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type MakeRequestUsecase interface {
	Request(context.CancelFunc, *entity.Request) *entity.Response
}
