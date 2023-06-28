package middleware

import (
	"errors"
	"github.com/yonisaka/go-boilerplate/internal/consts"
	"net/http"
	"strings"
)

func ValidatesContentType(r *http.Request) error {
	if ct := strings.ToLower(r.Header.Get(consts.HeaderContentTypeKey)); ct != consts.HeaderContentTypeJSON {
		return errors.New("invalid content-type")
	}

	return nil
}
