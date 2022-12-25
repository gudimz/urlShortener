package main

import (
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
	)
	logger.Infoln("Create router")
	router := httprouter.New()

	handler := shorten.NewHandler(logger)
	handler.Register(router)

	run(router)
}

func run(router *httprouter.Router) {
	var (
		logger = logging.GetLogger()
	)
	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		logger.Fatalln(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infoln("Shorten listening port 127.0.0.0:9000")
	logger.Fatalln(server.Serve(listener))
}
