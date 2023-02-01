package service

import (
	"context"

	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type MakeHttpRequestService struct {
}

func (MakeHttpRequestService) Request(cancelCtx context.CancelFunc, req *entity.Request, headers map[string]string) *entity.Response {
	return common.MakeLightweightRequest(cancelCtx, req, headers)
}
