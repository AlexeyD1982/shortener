package handler

import (
	"io"
	"log"
	"net/http"
)

var url string

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong request", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func SetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong request", http.StatusBadRequest)
		return
	}

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

	url = string(content)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("http://" + r.Host + "/EwHXdJfB"))

	if err != nil {
		log.Fatal(err)
	}
}
