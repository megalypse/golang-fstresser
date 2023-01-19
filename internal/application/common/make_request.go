package common

import (
	"bytes"
	"net/http"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

func MakeLightweightRequest[T entity.DataHolder](method string, req *entity.Request) entity.Response {
	client := http.Client{}

	httpRequest, _ := http.NewRequest(
		method,
		req.Url,
		bytes.NewBuffer(req.BytesBody),
	)

	for k, v := range req.Headers {
		httpRequest.Header.Add(k, v)
	}

	// TODO: handle error
	client.Do(httpRequest)

	return entity.Response{
		StatusCode: 0,
	}
}
