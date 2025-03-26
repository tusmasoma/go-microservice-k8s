package request

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Param(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
