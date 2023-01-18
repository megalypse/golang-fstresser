package usecase

import "github.com/megalypse/golang-fstresser/internal/domain/entity"

type PostFormUsecase interface {
	PostForm(entity.PostformRequest) *entity.Response
}
