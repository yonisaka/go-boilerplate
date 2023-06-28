package http

import (
	"context"
	"encoding/json"
	"github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"
	"github.com/yonisaka/go-boilerplate/internal/consts"
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/internal/http/middleware"
	"github.com/yonisaka/go-boilerplate/pkg/locales"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
	"net/http"
	"runtime/debug"
)

func (r *router) handle(hfn httpHandlerFunc, handler httphandler.Handler, mdws ...middleware.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := dto.HTTPResponse{
					Code:    http.StatusInternalServerError,
					Message: `Something went wrong, please try again later`,
				}

				res.GenerateMessage()
				logger.Error(logger.MessageFormat("error %v", string(debug.Stack())))
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
			r.response(w, dto.HTTPResponse{
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

func (r *router) response(w http.ResponseWriter, resp dto.HTTPResponse) {
	if resp.IsStream() {
		w.Header().Set(consts.HeaderContentTypeKey, resp.ContentType())

		defer func() {
			resp.GenerateMessage()
			w.WriteHeader(resp.GetCode())
			_, _ = w.Write(resp.DataStream())
		}()
	}

	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)

	defer func() {
		resp.GenerateMessage()
		w.WriteHeader(resp.GetCode())
		_ = json.NewEncoder(w).Encode(resp)
	}()
}

func (r *router) defaultLang(lang string) string {
	if len(lang) == 0 {
		return r.cfg.App.DefaultLang
	}

	return lang
}
