package main

import (
	"io"
	"log"
	"net/http"
)

var url string

func postEndpoint(w http.ResponseWriter, r *http.Request) {
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
	_, err = w.Write([]byte("http://localhost:8080/EwHXdJfB"))

	if err != nil {
		log.Fatal(err)
	}
}

func getEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong request", http.StatusBadRequest)
		return
	}

	w.Header().Set("location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", postEndpoint)
	mux.HandleFunc("/EwHXdJfB", getEndpoint)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
