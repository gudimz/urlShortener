package main

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/gudimz/urlShortener/pkg/logger"
	"github.com/gudimz/urlShortener/pkg/postgres"

	shortenerRepo "github.com/gudimz/urlShortener/internal/app/repository/psql/shortener"
	"github.com/gudimz/urlShortener/internal/app/service"
	"github.com/gudimz/urlShortener/internal/app/transport/rest/server"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
)

func main() {
	var logConfig logger.Config
	logConfig.ParseConfigFromEnv()
	log := logger.New(logConfig)

	log.Info("Trying to read config file")
	cfg := ds.GetConfig()

	log.Info("Trying to connect to db...")
	dbPool, err := postgres.NewClient(context.Background(), cfg.Postgres)
	if err != nil {
		log.Error("Failed to connect to db", zap.Error(err))
	}
	defer dbPool.Close()

	var (
		repository = shortenerRepo.NewRepository(dbPool, log)
		shortener  = service.New(repository)
		srv        = server.New(shortener, log)
	)

	run(log, srv, cfg)
}

func run(log *logger.Log, srv *server.Server, cfg *ds.Config) {
	log.Info("Shorten listening port", zap.String("port", cfg.Server.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), srv)
	if err != nil {
		log.Fatal("error running server", zap.Error(err))
	}
}
