package di

import (
	"net/http"

	"github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/pkg/routerkit"
)

// httpHandlerFunc is a contract http handler for router
type httpHandlerFunc func(request *http.Request, handler httphandler.Handler) dto.HTTPResponse

// Router is a contract router and must implement this interface
type Router interface {
	Route() *routerkit.Router
}
