package entity

type Request interface {
	GetUrl() string
	ToBytesBody() []byte
	GetHeaders() map[string]string
}

type PostformRequest interface {
	GetUrl() string
	ToMapBody() map[string][]string
	GetHeaders() map[string]string
}
