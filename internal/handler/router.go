package handler

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/AlexeyD1982/shortener/internal/config"
	localErrors "github.com/AlexeyD1982/shortener/pkg/errors"
	"github.com/go-chi/chi/v5"
)

type URLService interface {
	SaveURL(originURL string) string
	ResolveURL(shortURL string) (string, error)
}

type URLHandler struct {
	urlService URLService
	cfg        *config.Conf
}

func NewURLHandler(urlService URLService, cfg *config.Conf) *URLHandler {
	return &URLHandler{urlService: urlService, cfg: cfg}
}

func (h *URLHandler) InitRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.handlePost())
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

func ErrorHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}
