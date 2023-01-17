package usecasesimpl

import (
	"github.com/megalypse/golang-fstresser/internal/application/common"
	"github.com/megalypse/golang-fstresser/internal/domain/entity"
)

type HttpPost struct{}

func (HttpPost) Post(req entity.Request) entity.Response[entity.Void] {
	return common.MakeLightweightRequest[entity.Void]("POST", req)
}
