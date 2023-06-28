package dto

import (
	"github.com/yonisaka/go-boilerplate/config"
	"net/http"
)

type HttpData struct {
	Request     *http.Request
	Config      *config.Config
	ServiceType string
	BytesValue  []byte
}
