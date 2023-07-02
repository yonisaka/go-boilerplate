package httphandler

import (
	"net/http"

	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/internal/usecases"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

// healthHandler is a struct for health handler
type healthHandler struct {
	healthUsecase usecases.HealthUsecase
}

// NewHealthHandler is a constructor function for health handler
func NewHealthHandler(
	healthUsecase usecases.HealthUsecase,
) Handler {
	return &healthHandler{
		healthUsecase: healthUsecase,
	}
}

// Handle is a function to handle health check
func (h *healthHandler) Handle(req *http.Request) dto.HTTPResponse {
	logger.Info(logger.MessageFormat("health.check"), logger.EventName("health.check"))

	message, err := h.healthUsecase.Liveness(req.Context())
	if err != nil {
		return *dto.NewResponse().
			WithCode(http.StatusInternalServerError).
			WithMessage("Internal Server Error")
	}

	return *dto.NewResponse().WithCode(http.StatusOK).WithMessage(message)
}
