package repository

//go:generate rm -f ./health_mock.go
//go:generate mockgen -destination health_mock.go -package repository -mock_names HealthRepo=GoMockHealthRepo -source health.go

import "context"

// HealthRepo is a health repository interface.
type HealthRepo interface {
	GetLiveness(ctx context.Context) error
}
