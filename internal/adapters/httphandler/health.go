package httphandler

import (
	"net/http"

	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

type healthHandler struct {
}

func NewHealthHandler() Handler {
	return &healthHandler{}
}

func (h *healthHandler) Handle(req *http.Request) dto.HTTPResponse {
	logger.Info(logger.MessageFormat("health.check"), logger.EventName("health.check"))

	return *dto.NewResponse().WithCode(http.StatusOK).WithMessage("OK")
}
