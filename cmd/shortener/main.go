package main

import (
	"log"
	"net/http"

	"github.com/AlexeyD1982/shortener/internal/config"
	"github.com/AlexeyD1982/shortener/internal/handler"
	"github.com/AlexeyD1982/shortener/internal/repository/inmemory"
	"github.com/AlexeyD1982/shortener/internal/service"
)

func main() {
	cfg := config.Get(true)
	err := config.InitLogger("info")
	if err != nil {
		log.Fatal(err)
	}

	repo := inmemory.NewStorage()
	srv := service.NewURLService(repo)

	r := handler.NewURLHandler(srv, cfg)
	log.Fatal(http.ListenAndServe(cfg.Host, r.InitRouter()))
}
