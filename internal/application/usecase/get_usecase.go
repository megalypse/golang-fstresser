package usecasesimpl

import (
	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type HttpGet struct{}

func (HttpGet) Get(req entity.Request) entity.Response[entity.Void] {
	return common.MakeLightweightRequest[entity.Void]("GET", req)
}
