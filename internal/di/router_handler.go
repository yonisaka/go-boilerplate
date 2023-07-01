package di

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"
	"github.com/yonisaka/go-boilerplate/internal/consts"
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/internal/middleware"
	"github.com/yonisaka/go-boilerplate/pkg/locales"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
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

				func() {
					res.GenerateMessage()
					w.WriteHeader(res.GetCode())
					_ = json.NewEncoder(w).Encode(res)
				}()

				return
			}
		}()

		ctx := context.WithValue(req.Context(), consts.CtxAccess, map[string]interface{}{
			"path":      req.URL.Path,
			"remote_ip": req.RemoteAddr,
			"method":    req.Method,
		})

		lang := r.defaultLang(req.Header.Get(consts.HeaderAcceptLanguage))
		lang = locales.ParseIOSLang(lang, r.cfg.App.DefaultLang)
		ctx = locales.WithAcceptLanguage(ctx, lang)
		req = req.WithContext(ctx)

		if err := middleware.FilterFunc(req, mdws); err != nil {
			r.response(w, dto.HTTPResponse{
				Lang:   lang,
				Errors: err,
			})

			return
		}

		resp := hfn(req, handler)
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
	} else {
		w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)

		defer func() {
			func() {
				resp.GenerateMessage()
				w.WriteHeader(resp.GetCode())
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					logger.Error(logger.MessageFormat("error %v", err))
				}
			}()
		}()
	}
}

func (r *router) defaultLang(lang string) string {
	if len(lang) == 0 {
		return r.cfg.App.DefaultLang
	}

	return lang
}
