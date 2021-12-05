package middleware

import (
	"net/http"

	"github.com/practicalgo/book-exercise-solutions/chap6/exercise3/config"
)

func RegisterMiddleware(mux *http.ServeMux, c config.AppConfig) http.Handler {
	return loggingMiddleware(panicMiddleware(mux, c), c)
}
