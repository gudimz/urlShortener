package server

import (
	"github.com/gudimz/urlShortener/internal/shorten"
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Server struct {
	e       *echo.Echo
	logger  *logging.Logger
	shorten *shorten.Service
}

func NewServer(shorten *shorten.Service, logger *logging.Logger) *Server {
	srv := &Server{
		e:       echo.New(),
		logger:  logger,
		shorten: shorten,
	}
	srv.NewRouter()

	return srv

}

func (s *Server) NewRouter() {
	s.e.HideBanner = true
	s.e.Validator = NewValidator()

	s.e.Pre(middleware.RemoveTrailingSlash())
	s.e.Use(middleware.RequestID())

	s.RegisterRoutes()
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.e.ServeHTTP(writer, request)
}

func (s *Server) RegisterRoutes() {
	handler := NewHandler(s.shorten, s.logger)

	s.e.GET("/:short_url", handler.Redirect)
	g := s.e.Group("/api/v1")
	g.GET("/:short_url", handler.GetShorten)
	g.POST("/create", handler.CreateShorten)
	g.DELETE("/delete/:short_url", handler.DeleteShorten)
}
