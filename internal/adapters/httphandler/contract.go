package httphandler

import (
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"net/http"
)

type Handler interface {
	Handle(data *http.Request) dto.HTTPResponse
}
