package service

import (
	"log"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

func MakeRequest(req *entity.Request, requestService entity.RequestService) entity.Response {
	switch req.Method {
	case "GET":
		return requestService.Get(req)
	case "POST":
		return requestService.Post(req)
	case "POSTFORM":
		return requestService.PostForm(req)
	default:
		log.Fatalf("Http method not supported: %q", req.Method)
	}

	return entity.Response{}
}
