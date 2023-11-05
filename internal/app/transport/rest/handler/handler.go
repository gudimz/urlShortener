package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/samber/mo"
	"go.uber.org/zap"

	"github.com/gudimz/urlShortener/internal/app/service"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
	"github.com/gudimz/urlShortener/pkg/logger"
)

type Handler struct {
	log       *logger.Log
	shortener *service.Service
}

func New(shortener *service.Service, log *logger.Log) *Handler {
	return &Handler{
		log:       log,
		shortener: shortener,
	}
}

type request struct {
	Url      string `json:"url" validate:"required,url"`
	ShortUrl string `json:"short_url,omitempty" validate:"omitempty,alphanum"`
}

type response struct {
	Shortener *ds.Shorten `json:"shortener,omitempty"`
}

func (h *Handler) CreateShorten(ctx echo.Context) error {
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

	input := ds.InputShorten{
		ShortenUrl: shortenUrl,
		OriginUrl:  req.Url,
	}

	h.log.Info("create shorten for short url", zap.String("shortenUrl", input.ShortenUrl.OrEmpty()))
	shortener, err := h.shortener.CreateShorten(ctx.Request().Context(), input)
	if err != nil {
		if err == ds.ErrShortURLAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, "short url already exist")
		}

		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("can't create short url \"%v\"", input.ShortenUrl))
	}

	return ctx.JSON(http.StatusOK, response{Shortener: shortener})
}

func (h *Handler) Redirect(ctx echo.Context) error {
	shortUrl := ctx.Param("short_url")
	h.log.With(zap.String("shortUrl", shortUrl))
	h.log.Info("redirect for short url")
	originUrl, err := h.shortener.Redirect(ctx.Request().Context(), shortUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("short url %q not found", shortUrl))
		}

		h.log.Error("error getting redirect for short url", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("can't get url by short url %q", shortUrl))
	}
	return ctx.Redirect(http.StatusMovedPermanently, originUrl)
}

func (h *Handler) GetShorten(ctx echo.Context) error {
	shortUrl := ctx.Param("short_url")
	h.log.With(zap.String("shortUrl", shortUrl))
	h.log.Info("get shorten from db for short url")
	shortenInfo, err := h.shortener.GetShorten(ctx.Request().Context(), shortUrl)
	if err != nil {
		if err == ds.ErrShortUrlNotFound {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("short url %q not found", shortUrl))
		}

		h.log.Error("error getting GetShorten for short url", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to get shorten for  %q", shortUrl))
	}
	return ctx.JSON(http.StatusOK, shortenInfo)
}

func (h *Handler) DeleteShorten(ctx echo.Context) error {
	shortUrl := ctx.Param("short_url")
	h.log.With(zap.String("shortUrl", shortUrl))
	h.log.Info("delete shorten from db for short url")
	err := h.shortener.DeleteShorten(ctx.Request().Context(), shortUrl)
	if err != nil {
		if err == ds.ErrShortUrlNotFound {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("short url %q not found", shortUrl))
		}
		h.log.Error("error deleting GetShorten for short url", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("failed to delete shorten for  %q", shortUrl))
	}

	return ctx.NoContent(http.StatusNoContent)
}
