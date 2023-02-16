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

/*
Generic bad response to represent failed requests and avoid
creation of multiple unnecessary instances.

A downside is that with the current implementation it's not possible
to track down different http status codes than 500 from the response.
*/
var badResponse entity.Response

/*
Generic success response to represent succeeded requests and avoid
creation of multiple unnecessary instances.

A downside is that with the current implementation it's not possible
to track down different http status codes than 200 from the response.
*/
var successResponse entity.Response

/*
HttpClient pool created to reduce GC overhead.
*/
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

/*
MakeLightweightRequest is meant to use as little resources from the the machine hosting the stresser as possible.
*/
func MakeLightweightRequest(ctx context.Context, cancelCtx context.CancelFunc, req *entity.Request, headers map[string]string) *entity.Response {
	client := clientPool.Get().(*http.Client)
	defer clientPool.Put(client)

	httpRequest, err := http.NewRequest(
		req.Method,
		req.Url,
		bytes.NewBuffer(req.BytesBody),
	)
	httpRequest.Close = true

	if err != nil {
		GetLogger(ctx).Log(err.Error())
		cancelCtx()
		return &badResponse
	}

	// Global headers are added here
	for k, v := range headers {
		httpRequest.Header.Add(k, v)

	}

	// and then, the request specific headers are added here,
	// overwriting global header if the keys overlap.
	for k, v := range req.Headers {
		httpRequest.Header.Add(k, v)
	}

	res, err := client.Do(httpRequest)
	if err != nil {
		GetLogger(ctx).SilentLog(err.Error())
		return &badResponse
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			GetLogger(ctx).Log(err.Error())
			cancelCtx()
			return &badResponse
		}

		// Here we are silently logging the failed request. Silently for the enduser console does not get spammed with tons
		// of failed messages. It can be checked later on the created log file.
		GetLogger(ctx).SilentLog((fmt.Sprintf("Request failed with status code %d. Url: %q, Body: %q", res.StatusCode, req.Url, string(bytes))))
		return &badResponse
	}

	return &successResponse
}

// makeCustomDialer to make easier the creation of a custom dialer
func makeCustomDialer(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
