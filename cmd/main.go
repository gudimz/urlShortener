package main

import (
	"github.com/gudimz/urlShortener/internal/shorten"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Println("Create router")
	router := httprouter.New()

	handler := shorten.NewHandler()
	handler.Register(router)

	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Shorten listening port 127.0.0.0:9000")
	log.Fatalln(server.Serve(listener))
}
