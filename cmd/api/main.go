package main

import (
	"context"
	"fmt"
	"net/http"

	shortenRepo "github.com/gudimz/urlShortener/internal/app/repository/psql/shorten"
	"github.com/gudimz/urlShortener/internal/app/service"
	"github.com/gudimz/urlShortener/internal/app/transport/rest/server"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/gudimz/urlShortener/pkg/postgres"
)

func main() {

	var (
		logger = logging.GetLogger()
		cfg    = ds.GetConfig()
	)

	logger.Infoln("Trying to connect to db...")
	dbPool, err := postgres.NewClient(context.Background(), cfg.Postgres)
	if err != nil {
		logger.Fatalln(err)
	}
	defer dbPool.Close()

	var (
		repository = shortenRepo.NewRepository(dbPool, logger)
		shortener  = service.NewService(repository)
		srv        = server.New(shortener, logger)
	)

	run(srv, cfg)
}

func run(srv *server.Server, cfg *ds.Config) {
	var (
		logger = logging.GetLogger()
	)
	logger.Infoln(fmt.Sprintf("Shorten listening port :%s", cfg.Server.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), srv)
	if err != nil {
		logger.Fatalf("error running server: %v", err)
	}
}
