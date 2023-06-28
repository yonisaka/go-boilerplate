package http

import (
	"context"
	"github.com/yonisaka/go-boilerplate/config"
	"github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"
	"github.com/yonisaka/go-boilerplate/internal/consts"
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/pkg/msg"
	"net/http"
)

func httpRequest(request *http.Request, handler httphandler.Handler, conf *config.Config) dto.HTTPResponse {
	if !msg.GetAvailableLang(200, request.Header.Get(consts.HeaderLanguageKey)) {
		request.Header.Set(consts.HeaderLanguageKey, conf.App.DefaultLang)
	}

	ctx := context.WithValue(request.Context(), consts.CtxLang, request.Header.Get(consts.HeaderLanguageKey))

	req := request.WithContext(ctx)

	data := &dto.HTTPData{
		Request:     req,
		Config:      conf,
		ServiceType: consts.ServiceTypeHTTP,
	}

	return handler.Handle(data)
}
