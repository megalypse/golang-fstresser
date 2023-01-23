package entity

import "context"

type RequestService interface {
	Get(context.CancelFunc, *Request) *Response
	Post(context.CancelFunc, *Request) *Response
	PostForm(context.CancelFunc, *Request) *Response
}
