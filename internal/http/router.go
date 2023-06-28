package http

import (
	"github.com/yonisaka/go-boilerplate/config"
	"github.com/yonisaka/go-boilerplate/internal/di"
	"github.com/yonisaka/go-boilerplate/pkg/routerkit"
	"net/http"
)

type router struct {
	cfg    *config.Config
	router *routerkit.Router
}

func NewRouter() Router {
	cfg := di.GetConfig()
	return &router{cfg: cfg, router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.Name))}
}

func (r *router) Route() *routerkit.Router {
	root := r.router.PathPrefix("/").Subrouter()

	healthHandler := di.GetHealthHandler()

	root.HandleFunc("/liveness", r.handle(
		httpRequest,
		healthHandler,
	)).Methods(http.MethodGet)

	return r.router
}
