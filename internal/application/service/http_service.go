package service

import (
	"io"
	"log"
	"net/http"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type HttpService struct{}

func (HttpService) Get(req *entity.Request) entity.Response {
	return common.MakeLightweightRequest[entity.Void]("GET", req)
}

func (HttpService) Post(req *entity.Request) entity.Response {
	return common.MakeLightweightRequest[entity.Void]("POST", req)
}

func (HttpService) PostForm(req *entity.Request) entity.Response {
	res, err := http.PostForm(req.Url, req.MapBody)
	if err != nil {
		log.Fatal(err.Error())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	return entity.Response{
		StatusCode: res.StatusCode,
		Body:       body,
	}
}
