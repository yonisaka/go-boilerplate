package httphandler

import (
	"github.com/yonisaka/go-boilerplate/internal/dto"
)

type Handler interface {
	Handle(data *dto.HttpData) dto.HttpResponse
}
