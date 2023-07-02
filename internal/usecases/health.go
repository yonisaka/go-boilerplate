package usecases

import "context"

// Liveness is a function to check liveness
func (u *healthUsecase) Liveness(ctx context.Context) (string, error) {
	return "OK", nil
}
