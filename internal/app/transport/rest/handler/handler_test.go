package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	shortenRepo "github.com/gudimz/urlShortener/internal/app/repository/psql/shorten"
	"github.com/gudimz/urlShortener/internal/app/service"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
	"github.com/gudimz/urlShortener/pkg/logger"
	"github.com/gudimz/urlShortener/pkg/postgres"
)

var resp struct {
	Message string `json:"message"`
}

func TestHandler(t *testing.T) {
	// save current path
	testDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//change directory for read config
	if err := os.Chdir("../../../../.."); err != nil {
		panic(err)
	}

	var (
		cfg = ds.GetConfig()
	)

	log := logger.New(logger.Config{
		LogLevel: "debug",
	})

	dbPool, err := postgres.NewClient(context.Background(), cfg.Postgres)
	if err != nil {
		t.Fatal(err)
	}
	defer dbPool.Close()

	var (
		repository = shortenRepo.NewRepository(dbPool, log)
		shortener  = service.NewService(repository)
	)

	t.Run("Create new short url for a given URL", func(t *testing.T) {
		const body = `{"short_url": "youtube","url": "https://www.youtube.com"}`
		var (
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
			e        = echo.New()
			ctx      = e.NewContext(request, recorder)
			handler  = New(shortener, log)
		)
		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.NoError(t, handler.CreateShorten(ctx))
		assert.Equal(t, http.StatusOK, recorder.Code)

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp), &resp)
		assert.NotEmpty(t, resp.Message)
	})

	t.Run("Check two identical short_url", func(t *testing.T) {
		const body = `{"short_url": "google","url": "https://www.google.com"}`
		var (
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
			e        = echo.New()
			ctx      = e.NewContext(request, recorder)
			handler  = New(shortener, log)
		)
		e.Validator = NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.NoError(t, handler.CreateShorten(ctx))
		assert.Equal(t, http.StatusOK, recorder.Code)

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp), &resp)
		assert.NotEmpty(t, resp.Message)

		recorder = httptest.NewRecorder()
		request = httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx = e.NewContext(request, recorder)

		var httpErr *echo.HTTPError
		require.ErrorAs(t, handler.CreateShorten(ctx), &httpErr)
		assert.Equal(t, http.StatusConflict, httpErr.Code)
		assert.Contains(t, httpErr.Message, "short url already exist")
	})

	t.Run("Check redirect success", func(t *testing.T) {
		const (
			shortUrl  = "google"
			originUrl = "https://www.google.com"
		)
		var (
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest(http.MethodGet, "/"+shortUrl, nil)
			e        = echo.New()
			ctx      = e.NewContext(request, recorder)
			handler  = New(shortener, log)
		)

		ctx.SetPath("/:short_url")
		ctx.SetParamNames("short_url")
		ctx.SetParamValues(shortUrl)

		require.NoError(t, handler.Redirect(ctx))
		assert.Equal(t, http.StatusMovedPermanently, recorder.Code)
		assert.Equal(t, originUrl, recorder.Header().Get("Location"))
	})

	t.Run("Check redirect not found", func(t *testing.T) {
		const (
			shortUrl = "not_google"
		)
		var (
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest(http.MethodGet, "/"+shortUrl, nil)
			e        = echo.New()
			ctx      = e.NewContext(request, recorder)
			handler  = New(shortener, log)
		)

		ctx.SetPath("/:short_url")
		ctx.SetParamNames("short_url")
		ctx.SetParamValues(shortUrl)

		var httpErr *echo.HTTPError
		require.ErrorAs(t, handler.Redirect(ctx), &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})

	t.Run("Delete all short url", func(t *testing.T) {
		const (
			shortUrlFirst  = "google"
			shortUrlSecond = "youtube"
		)
		var (
			recorder = httptest.NewRecorder()
			request  = httptest.NewRequest(http.MethodDelete, "/delete/", nil)
			e        = echo.New()
			ctx      = e.NewContext(request, recorder)
			handler  = New(shortener, log)
		)
		e.Validator = NewValidator()
		ctx.SetPath("/delete/:short_url")
		ctx.SetParamNames("short_url")
		ctx.SetParamValues(shortUrlFirst)

		require.NoError(t, handler.DeleteShorten(ctx))
		assert.Equal(t, http.StatusNoContent, recorder.Code)

		recorder = httptest.NewRecorder()
		request = httptest.NewRequest(http.MethodDelete, "/delete/", nil)
		ctx = e.NewContext(request, recorder)

		ctx.SetPath("/delete/:short_url")
		ctx.SetParamNames("short_url")
		ctx.SetParamValues(shortUrlSecond)

		require.NoError(t, handler.DeleteShorten(ctx))
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	// delete dir with logs
	err = os.RemoveAll(filepath.Join(testDir, "logs"))
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}
