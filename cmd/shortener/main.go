package main

import (
	"log"
	"net/http"

	"github.com/AlexeyD1982/shortener/internal/config"
	"github.com/AlexeyD1982/shortener/internal/handler"
)

func main() {
	cfg := config.Get()

	r := handler.URLRouter(make(map[string]string))
	log.Fatal(http.ListenAndServe(cfg.Host, r))
}
