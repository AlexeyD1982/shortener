package middleware

import (
	"net/http"

	"github.com/AlexeyD1982/shortener/internal/config"
	"go.uber.org/zap"
)

func RequestLogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Info(
			"got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		h.ServeHTTP(w, r)
	})
}
