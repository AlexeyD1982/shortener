package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLogMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantMethod string
		wantPath   string
	}{
		{
			name:       "GET request",
			method:     http.MethodGet,
			path:       "/api/shorten",
			wantMethod: http.MethodGet,
			wantPath:   "/api/shorten",
		},
		{
			name:       "POST request",
			method:     http.MethodPost,
			path:       "/",
			wantMethod: http.MethodPost,
			wantPath:   "/",
		},
		{
			name:       "request with query params",
			method:     http.MethodGet,
			path:       "/api/user?id=123",
			wantMethod: http.MethodGet,
			wantPath:   "/api/user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, observed := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := RequestLogMiddleware(logger)
			wrappedHandler := middleware(handler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rec, req)

			logs := observed.All()
			if len(logs) != 1 {
				t.Fatalf("expected 1 log entry, got %d", len(logs))
			}

			log := logs[0]
			if log.Message != "HTTP request" {
				t.Errorf("expected message 'HTTP request', got %q", log.Message)
			}

			methodField := log.ContextMap()["method"]
			if methodField != tt.wantMethod {
				t.Errorf("expected method %q, got %q", tt.wantMethod, methodField)
			}

			pathField := log.ContextMap()["path"]
			if pathField != tt.wantPath {
				t.Errorf("expected path %q, got %q", tt.wantPath, pathField)
			}
		})
	}
}

func TestResponseLogMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		wantStatusCode string
		wantSize       string
	}{
		{
			name:           "200 OK with body",
			statusCode:     http.StatusOK,
			responseBody:   "test response",
			wantStatusCode: "200",
			wantSize:       "13",
		},
		{
			name:           "201 Created",
			statusCode:     http.StatusCreated,
			responseBody:   "created",
			wantStatusCode: "201",
			wantSize:       "7",
		},
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			responseBody:   "not found",
			wantStatusCode: "404",
			wantSize:       "9",
		},
		{
			name:           "empty body",
			statusCode:     http.StatusNoContent,
			responseBody:   "",
			wantStatusCode: "204",
			wantSize:       "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, observed := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			})

			middleware := ResponseLogMiddleware(logger)
			wrappedHandler := middleware(nextHandler)
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			logs := observed.All()
			if len(logs) != 1 {
				t.Fatalf("expected 1 log entry, got %d", len(logs))
			}

			log := logs[0]
			if log.Message != "HTTP response" {
				t.Errorf("expected message 'HTTP response', got %q", log.Message)
			}

			statusField := log.ContextMap()["statusCode"]
			if statusField != tt.wantStatusCode {
				t.Errorf("expected statusCode %q, got %q", tt.wantStatusCode, statusField)
			}

			sizeField := log.ContextMap()["size"]
			if sizeField != tt.wantSize {
				t.Errorf("expected size %q, got %q", tt.wantSize, sizeField)
			}

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}

			if w.Body.String() != tt.responseBody {
				t.Errorf("expected body %q, got %q", tt.responseBody, w.Body.String())
			}
		})
	}
}
