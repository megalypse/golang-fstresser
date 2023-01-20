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

	res, err := client.Do(httpRequest)

	if err != nil {
		return entity.Response{
			StatusCode: 500,
			Mesaage:    err.Error(),
		}
	}

	return entity.Response{
		StatusCode: res.StatusCode,
	}
}
