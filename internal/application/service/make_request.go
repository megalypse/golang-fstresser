package service

import (
	"context"
	"log"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

func MakeRequest(closeCtx context.CancelFunc, req *entity.Request, requestService entity.RequestService) *entity.Response {
	switch req.Method {
	case "GET":
		return requestService.Get(closeCtx, req)
	case "POST":
		return requestService.Post(closeCtx, req)
	case "POSTFORM":
		return requestService.PostForm(closeCtx, req)
	default:
		log.Fatalf("Http method not supported: %q", req.Method)
	}

	return &entity.Response{}
}
