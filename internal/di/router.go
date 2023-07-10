package di

import (
	"net/http"

	"github.com/yonisaka/go-boilerplate/config"
	"github.com/yonisaka/go-boilerplate/pkg/routerkit"
	"github.com/yonisaka/go-boilerplate/pkg/ws"
)

type router struct {
	cfg    *config.Config
	router *routerkit.Router
}

func NewRouter() Router {
	cfg := GetConfig()
	return &router{cfg: cfg, router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.Name))}
}

func (r *router) Route() *routerkit.Router {
	root := r.router.PathPrefix("/").Subrouter()

	healthHandler := GetHealthHandler()

	root.HandleFunc("/liveness", r.handle(
		httpGateway,
		healthHandler,
	)).Methods(http.MethodGet)

	hub := ws.NewHub(r.cfg)

	root.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleWebsocket(hub, w, r)
	}).Methods(http.MethodGet)

	return r.router
}
