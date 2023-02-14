package common

type ctxKey string

func GetCtxKey(rawKey string) ctxKey {
	return ctxKey(rawKey)
}
