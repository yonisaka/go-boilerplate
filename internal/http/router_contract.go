package http

import (
	"github.com/yonisaka/go-boilerplate/config"
	"github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/pkg/routerkit"
	"net/http"
)

// httpHandlerFunc is a contract http handler for router
type httpHandlerFunc func(request *http.Request, handler httphandler.Handler, cfg *config.Config) dto.HttpResponse

// Router is a contract router and must implement this interface
type Router interface {
	Route() *routerkit.Router
}
