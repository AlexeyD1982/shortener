package main

import (
	"log"
	"net/http"

	"github.com/AlexeyD1982/shortener/internal/handler"
)

func main() {
	r := handler.URLRouter(make(map[string]string))
	log.Fatal(http.ListenAndServe(":8080", r))
}
