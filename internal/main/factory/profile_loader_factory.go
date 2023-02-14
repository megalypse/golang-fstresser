package factory

import (
	"log"
	"os"

	"github.com/megalypse/golang-fstresser/internal/application/service"
	"github.com/megalypse/golang-fstresser/internal/domain/usecase"
)

var makeRequestUsecase usecase.MakeRequestUsecase

func init() {
	makeRequestUsecase = service.MakeHttpRequestService{}
}

func MakeLocalProfileLoader() usecase.ProfileLoader {
	path := os.Getenv("FSTRESSER_PROFILES_PATH")
	if path == "" {
		log.Fatal("Profiles path not defined")
	}

	return service.LocalProfileLoader{
		MakeRequestUsecase: makeRequestUsecase,
		ProfilesPath:       path,
	}
}
