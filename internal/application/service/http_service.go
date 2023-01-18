package service

import (
	"log"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
)

type HttpService struct {
	getUsecase      usecase.GetUsecase
	postUsecase     usecase.PostUsecase
	postformUsecase usecase.PostFormUsecase
}

func MakeRequest(req *entity.Request, httpService HttpService) *entity.Response {
	switch req.Method {
	case "GET":
		return httpService.Get(req)
	case "POST":
		return httpService.Post(req)
	default:
		log.Fatalf("Http method not suported: %q", req.Method)
		return &entity.Response{}
	}
}

func (hs HttpService) Get(req *entity.Request) *entity.Response {
	return hs.getUsecase.Get(req)
}

func (hs HttpService) Post(req *entity.Request) *entity.Response {
	return hs.postUsecase.Post(req)
}

func (hs HttpService) Postform(req entity.PostformRequest) *entity.Response {
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
