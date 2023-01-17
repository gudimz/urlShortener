package shorten

import (
	"context"
	"github.com/gudimz/urlShortener/internal/config"
	"github.com/gudimz/urlShortener/internal/db/postgres"
	"github.com/gudimz/urlShortener/internal/model"
	shorten2 "github.com/gudimz/urlShortener/internal/storage/shorten"
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"testing"
)

const pathToConf = "../../config.yml"

func Test_CreateShortenUrl(t *testing.T) {
	var (
		cfg             = config.GetConfig(pathToConf)
		logger          = logging.GetLogger()
		generateShorten = ""
		customShorten   = ""
	)

	dbPool, err := postgres.NewClient(context.Background(), cfg.Postgres)
	if err != nil {
		t.Fatal(err)
	}
	defer dbPool.Close()

	var (
		storage = shorten2.NewStorage(dbPool, logger)
		service = NewService(storage)
	)

	t.Run("Create new short url with GenerateShortenUrl()", func(t *testing.T) {
		var (
			inputShorten = model.InputShorten{
				OriginUrl: "https://youtube.com/",
			}
		)

		shorten, err := service.CreateShorten(context.Background(), inputShorten)
		generateShorten = shorten.ShortUrl
		assert.NoError(t, err)
		assert.Equal(t, shorten.OriginUrl, "https://youtube.com/")
		assert.NotZero(t, shorten.DateCreated)
		assert.NotZero(t, shorten.DateUpdated)
	})

	t.Run("Create new short url with custom short url", func(t *testing.T) {
		var (
			inputShorten = model.InputShorten{
				CustomUrl: mo.Some("google"),
				OriginUrl: "https://google.com/",
			}
		)

		shorten, err := service.CreateShorten(context.Background(), inputShorten)
		customShorten = shorten.ShortUrl
		assert.NoError(t, err)
		assert.Equal(t, shorten.ShortUrl, "google")
		assert.Equal(t, shorten.OriginUrl, "https://google.com/")
		assert.NotZero(t, shorten.DateCreated)
		assert.NotZero(t, shorten.DateUpdated)
	})

	t.Run("Delete short url random and custom", func(t *testing.T) {

		errGen := service.DeleteShorten(context.Background(), generateShorten)
		assert.NoError(t, errGen)

		errCustom := service.DeleteShorten(context.Background(), customShorten)
		assert.NoError(t, errCustom)
	})

}
