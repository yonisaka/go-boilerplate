package httphandler

import (
	"net/http"

	"github.com/yonisaka/go-boilerplate/internal/dto"
)

type Handler interface {
	Handle(data *http.Request) dto.HTTPResponse
}
