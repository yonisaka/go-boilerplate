package middleware

import (
	"github.com/yonisaka/go-boilerplate/config"
	"net/http"
)

// MiddlewareFunc is contract for middleware and must implement this type for http if need middleware http request
type MiddlewareFunc func(r *http.Request, conf *config.Config) error

// FilterFunc is a iterator resolver in each middleware registered
func FilterFunc(conf *config.Config, r *http.Request, mfs []MiddlewareFunc) error {
	for _, mf := range mfs {
		if err := mf(r, conf); err != nil {
			return err
		}
	}

	return nil
}
