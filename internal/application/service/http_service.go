package service

import (
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
)

type HttpService struct {
	getUsecase      usecase.GetUsecase
	postUsecase     usecase.PostUsecase
	postformUsecase usecase.PostFormUsecase
}

func (hs HttpService) Get(req entity.Request) entity.Response[entity.Void] {
	return hs.getUsecase.Get(req)
}

func (hs HttpService) Post(req entity.Request) entity.Response[entity.Void] {
	return hs.postUsecase.Post(req)
}

func (hs HttpService) Postform(req entity.PostformRequest) entity.Response[entity.Void] {
	return hs.postformUsecase.PostForm(req)
}

func MakeHttpService(
	getUsecase usecase.GetUsecase,
	postUsecase usecase.PostUsecase,
	postformUsecase usecase.PostFormUsecase,
) HttpService {
	return HttpService{
		getUsecase:      getUsecase,
		postUsecase:     postUsecase,
		postformUsecase: postformUsecase,
	}
}
