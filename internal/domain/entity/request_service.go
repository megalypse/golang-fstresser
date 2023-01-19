package entity

type RequestService interface {
	Get(req *Request) Response
	Post(req *Request) Response
	PostForm(req *Request) Response
}
