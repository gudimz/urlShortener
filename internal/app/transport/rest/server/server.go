package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gudimz/urlShortener/internal/app/service"
	"github.com/gudimz/urlShortener/internal/app/transport/rest/handler"
	"github.com/gudimz/urlShortener/pkg/logger"
)

type Server struct {
	e       *echo.Echo
	log     *logger.Log
	shorten *service.Service
}

func New(shorten *service.Service, log *logger.Log) *Server {
	srv := &Server{
		e:       echo.New(),
		log:     log,
		shorten: shorten,
	}
	srv.NewRouter()

	return srv
}

func (s *Server) NewRouter() {
	s.e.HideBanner = true
	s.e.Validator = handler.NewValidator()

	s.e.Pre(middleware.RemoveTrailingSlash())
	s.e.Use(middleware.RequestID())

	s.RegisterRoutes()
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.e.ServeHTTP(writer, request)
}

func (s *Server) RegisterRoutes() {
	h := handler.New(s.shorten, s.log)

	s.e.GET("/:short_url", h.Redirect)

	g := s.e.Group("/api/v1")
	g.GET("/:short_url", h.GetShorten)
	g.POST("/create", h.CreateShorten)
	g.DELETE("/delete/:short_url", h.DeleteShorten)
}
