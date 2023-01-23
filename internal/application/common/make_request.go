package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

var badResponse entity.Response
var successResponse entity.Response
var clientPool sync.Pool

func init() {
	badResponse = entity.Response{
		StatusCode: 500,
	}

	successResponse = entity.Response{
		StatusCode: 200,
	}

	clientPool = sync.Pool{
		New: func() any { return new(http.Client) },
	}
}

func MakeLightweightRequest(closeCtx context.CancelFunc, req *entity.Request) *entity.Response {
	client := clientPool.Get().(*http.Client)
	defer clientPool.Put(client)

	httpRequest, err := http.NewRequest(
		req.Method,
		req.Url,
		bytes.NewBuffer(req.BytesBody),
	)

	if err != nil {
		GetLogger().Log(err.Error())
		closeCtx()
	}

	for k, v := range req.Headers {
		httpRequest.Header.Add(k, v)
	}

	res, err := client.Do(httpRequest)

	if err != nil {
		GetLogger().Log(err.Error())
		closeCtx()

	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			GetLogger().Log(err.Error())
			closeCtx()
		}

		GetLogger().Log((fmt.Sprintf("Request failed with status code %d. Body: %q", res.StatusCode, string(bytes))))
		closeCtx()
	}

	return &successResponse
}
