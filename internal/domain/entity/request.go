package entity

type Request struct {
	Method    string
	Url       string
	BytesBody []byte
	Headers   map[string]string
}

type PostformRequest struct {
	Url     string
	MapBody map[string][]string
	Headers map[string]string
}
