package usecasesimpl

import (
	"io"
	"log"
	"net/http"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type HttpPostform struct{}

func (HttpPostform) PostForm(req entity.PostformRequest) (int, []byte) {
	res, err := http.PostForm(req.Url, req.MapBody)
	if err != nil {
		log.Fatal(err.Error())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	return res.StatusCode, body
}
