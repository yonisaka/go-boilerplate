package usecases_test

import (
	"github.com/yonisaka/go-boilerplate/internal/entities/repository"
	"github.com/yonisaka/go-boilerplate/internal/usecases"
)

type fields struct {
	healthRepo repository.HealthRepo
}

func sut(f fields) usecases.HealthUsecase {
	return usecases.NewHealthUsecase(
		f.healthRepo,
	)
}
