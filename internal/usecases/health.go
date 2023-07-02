package usecases

import "context"

// Liveness is a function to check liveness
func (u *healthUsecase) Liveness(ctx context.Context) (string, error) {
	err := u.healthRepo.GetLiveness(ctx)
	if err != nil {
		return "", err
	}

	return "OK", nil
}
