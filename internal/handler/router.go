package handler

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func URLHandler(urls map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("content-type") != "text/plain" {
				http.Error(w, "Wrong request", http.StatusBadRequest)
				return
			}

			if r.URL.Path != "/" {
				http.Error(w, "Wrong request", http.StatusBadRequest)
				return
			}

			content, err := io.ReadAll(r.Body)

			if err != nil {
				log.Fatal(err)
			}

			url := string(content)
			code := generateRandomString(8)

			if urls == nil {
				urls = make(map[string]string)
			}

			urls[code] = url
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write([]byte("http://" + r.Host + "/" + code))

			if err != nil {
				log.Fatal(err)
			}
			return
		} else if r.Method == http.MethodGet {
			path := r.URL.Path
			segments := strings.Split(path, "/")

			if len(segments) != 2 || segments[0] != "" || segments[1] == "" {
				http.Error(w, "Wrong request", http.StatusBadRequest)
				return
			}

			url, ok := urls[segments[1]]

			if !ok {
				http.Error(w, "Wrong request", http.StatusBadRequest)
			}

			w.Header().Set("location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		http.Error(w, "Wrong request", http.StatusBadRequest)
	}
}
