package main

import (
	"net/http"

	"github.com/AlexeyD1982/shortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.SetUrlHandler)
	mux.HandleFunc("/EwHXdJfB", handler.GetUrlHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
