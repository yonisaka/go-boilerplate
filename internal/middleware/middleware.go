package middleware

import (
	"net/http"
)

// Handle is contract for middleware and must implement this type for http if you need middleware http request
type Handle func(r *http.Request) error

// FilterFunc is a iterator resolver in each middleware registered
func FilterFunc(r *http.Request, mfs []Handle) error {
	for _, mf := range mfs {
		if err := mf(r); err != nil {
			return err
		}
	}

	return nil
}
