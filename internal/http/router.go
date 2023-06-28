package http

import (
	"context"
	"encoding/json"
	"github.com/yonisaka/go-boilerplate/config"
	"github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"
	"github.com/yonisaka/go-boilerplate/internal/consts"
	"github.com/yonisaka/go-boilerplate/internal/di"
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/internal/http/middleware"
	"github.com/yonisaka/go-boilerplate/pkg/locales"
	"github.com/yonisaka/go-boilerplate/pkg/routerkit"
	"net/http"
)

type router struct {
	cfg    *config.Config
	router *routerkit.Router
}

func NewRouter() Router {
	cfg := di.GetConfig()
	return &router{
		cfg:    cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.Name)),
	}
}

func (r *router) handle(hfn httpHandlerFunc, handler httphandler.Handler, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := dto.HttpResponse{
					Code:    http.StatusInternalServerError,
					Message: `Something went wrong, please try again later`,
				}

				res.GenerateMessage()
				//logger.Error(logger.MessageFormat("error %v", string(debug.Stack())))
				_ = json.NewEncoder(w).Encode(res)
				return
			}
		}()

		ctx := context.WithValue(req.Context(), "access", map[string]interface{}{
			"path":      req.URL.Path,
			"remote_ip": req.RemoteAddr,
			"method":    req.Method,
		})

		lang := r.defaultLang(req.Header.Get(consts.HeaderAcceptLanguage))
		lang = locales.ParseIOSLang(lang, r.cfg.App.DefaultLang)
		ctx = locales.WithAcceptLanguage(ctx, lang)
		req = req.WithContext(ctx)

		if err := middleware.FilterFunc(r.cfg, req, mdws); err != nil {
			r.response(w, dto.HttpResponse{
				Lang:   lang,
				Errors: err,
			})

			return
		}

		resp := hfn(req, handler, r.cfg)
		resp.Lang = lang
		r.response(w, resp)
	}
}

func (r *router) response(w http.ResponseWriter, resp dto.HttpResponse) {
	if resp.IsStream() {
		w.Header().Set(consts.HeaderContentTypeKey, resp.ContentType())

		defer func() {
			resp.GenerateMessage()
			w.WriteHeader(resp.GetCode())
			_, _ = w.Write(resp.DataStream())
		}()
		return
	}

	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)

	defer func() {
		resp.GenerateMessage()
		w.WriteHeader(resp.GetCode())
		_ = json.NewEncoder(w).Encode(resp)
	}()

	return
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

func (r *router) defaultLang(lang string) string {
	if len(lang) == 0 {
		return r.cfg.App.DefaultLang
	}

	return lang
}
