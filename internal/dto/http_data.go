package dto

import (
	"github.com/yonisaka/go-boilerplate/config"
	"net/http"
)

type HTTPData struct {
	Request     *http.Request
	Config      *config.Config
	ServiceType string
	BytesValue  []byte
}
