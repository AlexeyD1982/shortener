package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/AlexeyD1982/shortener/internal/config"
	"github.com/AlexeyD1982/shortener/internal/handler/middleware"
	"github.com/AlexeyD1982/shortener/internal/model"
	localErrors "github.com/AlexeyD1982/shortener/pkg/errors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type URLService interface {
	SaveURL(originURL string) string
	ResolveURL(shortURL string) (string, error)
}

type URLHandler struct {
	urlService URLService
	cfg        *config.Conf
	logger     *zap.Logger
}

func NewURLHandler(urlService URLService, cfg *config.Conf, logger *zap.Logger) *URLHandler {
	return &URLHandler{urlService: urlService, cfg: cfg, logger: logger}
}

func (h *URLHandler) InitRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogMiddleware(h.logger))
	r.Use(middleware.ResponseLogMiddleware(h.logger))
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.handlePost())
		r.Post("/api/shorten", h.handleShortenPost())
		r.Get("/{id}", h.handleGet())
		r.MethodNotAllowed(ErrorHandler)
		r.NotFound(ErrorHandler)
	})

	return r
}

func (h *URLHandler) handlePost() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "wrong request", http.StatusBadRequest)
			return
		}

		content, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		originURL := string(content)
		shortURL := h.urlService.SaveURL(originURL)

		w.WriteHeader(http.StatusCreated)
		resultPath, err := url.JoinPath(h.cfg.ResultHost, shortURL)

		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		_, err = w.Write([]byte(resultPath))

		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}
}

func (h *URLHandler) handleGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		originURL, err := h.urlService.ResolveURL(id)

		if err != nil {
			if errors.Is(err, localErrors.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		w.Header().Set("Location", originURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *URLHandler) handleShortenPost() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "wrong request", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		var req model.ApiShortenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.URL == "" {
			http.Error(w, "invalid request", http.StatusUnprocessableEntity)
			return
		}

		shortURL := h.urlService.SaveURL(req.URL)

		resultPath, err := url.JoinPath(h.cfg.ResultHost, shortURL)

		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		res := model.ApiShortenResponse{
			Result: resultPath,
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ErrorHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}
