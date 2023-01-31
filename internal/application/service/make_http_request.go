package service

import (
	"context"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type MakeHttpRequestService struct {
}

func (MakeHttpRequestService) Request(cancelCtx context.CancelFunc, req *entity.Request) *entity.Response {
	return common.MakeLightweightRequest(cancelCtx, req)
}
