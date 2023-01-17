package common

import (
	"bytes"
	"log"
	"net/http"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

func MakeLightweightRequest[T entity.DataHolder](method string, req entity.Request) entity.Response[T] {
	client := http.Client{}

	httpRequest, _ := http.NewRequest(
		method,
		req.Url,
		bytes.NewBuffer(req.BytesBody),
	)

	for k, v := range req.Headers {
		httpRequest.Header.Add(k, v)
	}

	res, err := client.Do(httpRequest)

	if err != nil {
		log.Fatal(err.Error())
	}

	return entity.Response[T]{
		StatusCode: res.StatusCode,
	}
}
