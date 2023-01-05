package main

import (
	"context"
	"fmt"
	"github.com/gudimz/urlShortener/internal/config"
	"github.com/gudimz/urlShortener/internal/db/postgres"
	"github.com/gudimz/urlShortener/internal/shorten"
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
)

func main() {
	var (
		logger = logging.GetLogger()
		cfg    = config.GetConfig()
		router = httprouter.New()
	)

	handler := shorten.NewHandler(logger)
	handler.Register(router)

	logger.Infoln("Trying to connect to db...")
	dbPool, err := postgres.NewClient(context.Background(), cfg.Postgres)
	if err != nil {
		logger.Fatalln(err)
	}
	defer dbPool.Close()

	run(router, cfg)
}

func run(router *httprouter.Router, cfg *config.Config) {
	var (
		logger = logging.GetLogger()
	)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.Ip, cfg.Listen.Port))
	if err != nil {
		logger.Fatalln(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infoln(fmt.Sprintf("Shorten listening port %s:%s", cfg.Listen.Ip, cfg.Listen.Port))
	logger.Fatalln(server.Serve(listener))
}
