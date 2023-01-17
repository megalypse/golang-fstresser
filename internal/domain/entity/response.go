package entity

type Response[T any] struct {
	Body       T
	StatusCode int
}
