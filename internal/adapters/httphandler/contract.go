package httphandler

import (
	"net/http"

	"github.com/yonisaka/go-boilerplate/internal/dto"
)

// Handler is an interface for http handler
type Handler interface {
	Handle(data *http.Request) dto.HTTPResponse
}
