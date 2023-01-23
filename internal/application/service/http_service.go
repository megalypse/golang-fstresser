package service

import (
	"context"
	"io"
	"net/http"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type HttpService struct{}

func (HttpService) Get(closeCtx context.CancelFunc, req *entity.Request) *entity.Response {
	return common.MakeLightweightRequest(closeCtx, req)
}

func (HttpService) Post(closeCtx context.CancelFunc, req *entity.Request) *entity.Response {
	return common.MakeLightweightRequest(closeCtx, req)
}

func (HttpService) PostForm(closeCtx context.CancelFunc, req *entity.Request) *entity.Response {
	res, err := http.PostForm(req.Url, req.MapBody)
	if err != nil {
		common.GetLogger().Log(err.Error())
		closeCtx()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		common.GetLogger().Log(err.Error())
		closeCtx()
	}

	return &entity.Response{
		StatusCode: res.StatusCode,
		Body:       body,
	}
}
