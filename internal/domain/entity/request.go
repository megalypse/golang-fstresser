package entity

type Request struct {
	Method    string
	Url       string
	BytesBody []byte
	MapBody   map[string][]string
	Headers   map[string]string
}
