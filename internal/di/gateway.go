package di

import (
	"context"
	"net/http"

	"github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"
	"github.com/yonisaka/go-boilerplate/internal/consts"
	"github.com/yonisaka/go-boilerplate/internal/dto"
	"github.com/yonisaka/go-boilerplate/pkg/msg"
)

func httpGateway(request *http.Request, handler httphandler.Handler) dto.HTTPResponse {
	cfg := GetConfig()
	if !msg.GetAvailableLang(200, request.Header.Get(consts.HeaderLanguageKey)) {
		request.Header.Set(consts.HeaderLanguageKey, cfg.App.DefaultLang)
	}

	ctx := context.WithValue(request.Context(), consts.CtxLang, request.Header.Get(consts.HeaderLanguageKey))

	req := request.WithContext(ctx)

	return handler.Handle(req)
}
