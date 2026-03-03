package main

import (
	"log"
	"net/http"

	"github.com/AlexeyD1982/shortener/internal/config"
	"github.com/AlexeyD1982/shortener/internal/handler"
)

func main() {
	config.ParseFlags()

	r := handler.URLRouter(make(map[string]string))
	log.Fatal(http.ListenAndServe(config.Conf.Host, r))
}
