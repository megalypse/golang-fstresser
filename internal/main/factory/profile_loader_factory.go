package factory

import (
	"github.com/megalypse/golang-fstresser/internal/application/service"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
)

var makeRequestUsecase usecase.MakeRequestUsecase

func init() {
	makeRequestUsecase = service.MakeHttpRequestService{}
}

func MakeLocalProfileLoader() usecase.ProfileLoader {
	return service.LocalProfileLoader{}
}
