package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

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
		New: func() any {
			customClient := http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: makeCustomDialer(&net.Dialer{
						KeepAlive: time.Nanosecond * -1,
					}),
					ForceAttemptHTTP2:     true,
					MaxIdleConns:          0,
					IdleConnTimeout:       5 * time.Second,
					TLSHandshakeTimeout:   0,
					ExpectContinueTimeout: 1 * time.Second,
				},
			}

			return &customClient
		},
	}
}

func MakeLightweightRequest(cancelCtx context.CancelFunc, req *entity.Request) *entity.Response {
	client := clientPool.Get().(*http.Client)
	defer clientPool.Put(client)

	httpRequest, err := http.NewRequest(
		req.Method,
		req.Url,
		bytes.NewBuffer(req.BytesBody),
	)
	httpRequest.Close = true

	if err != nil {
		GetLogger().Log(err.Error())
		cancelCtx()
		return &badResponse
	}

	for k, v := range req.Headers {
		httpRequest.Header.Add(k, v)
	}

	res, err := client.Do(httpRequest)
	if err != nil {
		GetLogger().SilentLog(err.Error())
		return &badResponse
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			GetLogger().Log(err.Error())
			cancelCtx()
			return &badResponse
		}

		GetLogger().SilentLog((fmt.Sprintf("Request failed with status code %d. Body: %q", res.StatusCode, string(bytes))))
		return &badResponse
	}

	return &successResponse
}

func makeCustomDialer(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
