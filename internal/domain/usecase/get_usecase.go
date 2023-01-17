package usecase

import "github.com/megalypse/golang-fstresser/internal/domain/entity"

type GetUsecase interface {
	Get(entity.Request) entity.Response[entity.Void]
}
