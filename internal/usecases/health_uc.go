package usecases

//go:generate rm -f ./health_uc_mock.go
//go:generate mockgen -destination health_uc_mock.go -package usecases -mock_names HealthUsecase=GoMockHealthUsecase -source health_uc.go

import (
	"context"

	"github.com/yonisaka/go-boilerplate/internal/entities/repository"
)

// HealthUsecase is an interface for health usecase
type HealthUsecase interface {
	Liveness(ctx context.Context) (string, error)
}

// Compile time implementation check
var _ HealthUsecase = (*healthUsecase)(nil)

// NewHealthUsecase is a constructor function for health usecase
func NewHealthUsecase(
	healthRepo repository.HealthRepo,
) HealthUsecase {
	return &healthUsecase{
		healthRepo: healthRepo,
	}
}

// healthUsecase is a struct for health usecase
type healthUsecase struct {
	healthRepo repository.HealthRepo
}
