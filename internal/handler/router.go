package handler

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func URLRouter(urls map[string]string) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", PostHandler(urls))
		r.Get("/{id}", GetHandler(urls))
		r.MethodNotAllowed(ErrorHandler)
		r.NotFound(ErrorHandler)
	})

	return r
}

func PostHandler(urls map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Wrong request", http.StatusBadRequest)
			return
		}

		content, err := io.ReadAll(r.Body)

		if err != nil {
			log.Fatal(err)
		}

		url := string(content)
		code := generateRandomString(8, time.Now().UnixNano())

		if urls == nil {
			urls = make(map[string]string)
		}

		urls[code] = url
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("http://" + r.Host + "/" + code))

		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetHandler(urls map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		url, ok := urls[id]

		if !ok {
			http.Error(w, "Wrong request", http.StatusBadRequest)
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func ErrorHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}
