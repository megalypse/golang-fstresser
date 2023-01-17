package service

import "github.com/megalypse/golang-fstresser/internal/domain/usecase"

type HttpService struct {
	GetUsecase      usecase.GetUsecase
	PostUsecase     usecase.PostUsecase
	PostformUsecase usecase.PostFormUsecase
}
