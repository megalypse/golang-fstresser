package common

// Type created to avoid using string type keys when
// getting or setting values to/from a context
type ctxKey string

// GetCtxKey provides a convenient way of using `ctxKey`.
// Use it whenever there's the need of getting/setting a value on a context.
func GetCtxKey(rawKey string) ctxKey {
	return ctxKey(rawKey)
}
