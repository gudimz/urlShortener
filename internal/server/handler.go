package server

import (
	"errors"
	"fmt"
	"github.com/gudimz/urlShortener/internal/config"
	"github.com/gudimz/urlShortener/internal/model"
	"github.com/gudimz/urlShortener/internal/shorten"
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
	"net/http"
	"strings"
)

type handler struct {
	logger    *logging.Logger
	shortener *shorten.Service
}

func NewHandler(shortener *shorten.Service, logger *logging.Logger) *handler {
	return &handler{
		logger:    logger,
		shortener: shortener,
	}
}

type request struct {
	Url      string `json:"url" validate:"required,url"`
	ShortUrl string `json:"short_url,omitempty" validate:"omitempty,alphanum"`
}

type response struct {
	Message string `json:"message,omitempty"`
}

func (h *handler) CreateShorten(ctx echo.Context) error {
	var req request
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}
	shortenUrl := mo.None[string]()
	if strings.TrimSpace(req.ShortUrl) != "" {
		shortenUrl = mo.Some(req.ShortUrl)
	}

	input := model.InputShorten{
		ShortenUrl: shortenUrl,
		OriginUrl:  req.Url,
	}

	h.logger.Infof("create shorten for short url \"%v\"", input.ShortenUrl)
	shortener, err := h.shortener.CreateShorten(ctx.Request().Context(), input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if strings.Compare(pgErr.Code, "23505") == 0 {
				return echo.NewHTTPError(http.StatusConflict, "short url already exist")
			}
		}
		h.logger.Errorf("error creating shorten: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("can't create short url \"%v\"", input.ShortenUrl))
	}
	message := fmt.Sprintf("%v://%v:%v/%v",
		config.GetConfig().Server.Protocol,
		config.GetConfig().Server.Host,
		config.GetConfig().Server.Port,
		shortener.ShortUrl,
	)
	return ctx.JSON(http.StatusOK, response{Message: message})
}

func (h *handler) Redirect(ctx echo.Context) error {
	shortUrl := ctx.Param("short_url")
	h.logger.Infof("redirect for short url %q", shortUrl)
	originUrl, err := h.shortener.Redirect(ctx.Request().Context(), shortUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("short url %q not found", shortUrl))
		}

		h.logger.Errorf("error getting redirect for short url %q: %v", shortUrl, err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("can't get url by short url %q", shortUrl))
	}
	return ctx.Redirect(http.StatusMovedPermanently, originUrl)
}

func (h *handler) GetShorten(ctx echo.Context) error {
	shortUrl := ctx.Param("short_url")
	h.logger.Infof("get shorten from db for short url %q", shortUrl)
	shortenInfo, err := h.shortener.GetShorten(ctx.Request().Context(), shortUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("short url %q not found", shortUrl))
		}

		h.logger.Errorf("error getting GetShorten for short url %q: %v", shortUrl, err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to get shorten for  %q", shortUrl))
	}
	return ctx.JSON(http.StatusOK, shortenInfo)
}

func (h *handler) DeleteShorten(ctx echo.Context) error {
	shortUrl := ctx.Param("short_url")
	h.logger.Infof("delete shorten from db for short url %q", shortUrl)
	count, err := h.shortener.DeleteShorten(ctx.Request().Context(), shortUrl)
	if count == 0 {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("short url %q not found", shortUrl))
	}
	if err != nil {
		h.logger.Errorf("error deleting GetShorten for short url %q: %v", shortUrl, err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to delete shorten for  %q", shortUrl))
	}
	return ctx.NoContent(http.StatusNoContent)
}
