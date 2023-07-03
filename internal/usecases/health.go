package usecases

import "context"

// Liveness is a function to check liveness
func (u *healthUsecase) Liveness(ctx context.Context) (string, error) {
	err := u.healthRepo.GetLiveness(ctx)
	if err != nil {
		return "", err
	}

	return "Life is not about waiting for the storms to pass; " +
		"it's about learning to dance in the rain and embracing every challenge as an opportunity for growth.", nil
}
