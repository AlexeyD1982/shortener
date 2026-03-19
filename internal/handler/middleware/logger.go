package middleware

import (
	"net/http"
	"strconv"

	"github.com/AlexeyD1982/shortener/internal/config"
	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.responseData.status = status
}

func RequestLogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Info(
			"HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		h.ServeHTTP(w, r)
	})
}

func ResponseLogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		config.Logger.Info(
			"HTTP response",
			zap.String("statusCode", strconv.Itoa(lw.responseData.status)),
			zap.String("size", strconv.Itoa(lw.responseData.size)),
		)
	})
}
