package main

import (
	"context"
	"fmt"
	"github.com/gudimz/urlShortener/internal/config"
	"github.com/gudimz/urlShortener/internal/db/postgres"
	"github.com/gudimz/urlShortener/internal/server"
	"github.com/gudimz/urlShortener/internal/shorten"
	shorten2 "github.com/gudimz/urlShortener/internal/storage/shorten"
	"github.com/gudimz/urlShortener/pkg/logging"
	"net/http"
)

func main() {

	var (
		logger = logging.GetLogger()
		cfg    = config.GetConfig()
	)

	logger.Infoln("Trying to connect to db...")
	dbPool, err := postgres.NewClient(context.Background(), cfg.Postgres)
	if err != nil {
		logger.Fatalln(err)
	}
	defer dbPool.Close()

	var (
		storage   = shorten2.NewStorage(dbPool, logger)
		shortener = shorten.NewService(storage)
		srv       = server.NewServer(shortener, logger)
	)

	run(srv, cfg)
}

func run(srv *server.Server, cfg *config.Config) {
	var (
		logger = logging.GetLogger()
	)
	logger.Infoln(fmt.Sprintf("Shorten listening port :%s", cfg.Server.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), srv)
	if err != nil {
		logger.Fatalf("error running server: %v", err)
	}
}
