package locales

import (
	"context"
	"strings"
)

type LangContextKey struct{}

func WithAcceptLanguage(ctx context.Context, lang string) context.Context {
	ctx = context.WithValue(ctx, LangContextKey{}, strings.ToLower(lang))
	return ctx
}

func LanguageFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(LangContextKey{}).(string)
	return v, ok
}
