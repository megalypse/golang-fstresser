package usecase

import "github.com/megalypse/golang-fstresser/internal/domain/entity"

type PostUsecase interface {
	Post(*entity.Request) *entity.Response
}
