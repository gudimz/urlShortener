package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"

	"github.com/gudimz/urlShortener/internal/app/repository/psql/shorten"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/gudimz/urlShortener/pkg/postgres"
)

func TestService(t *testing.T) {
	// save current path
	testDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//change directory for read config
	if err := os.Chdir("../../.."); err != nil {
		panic(err)
	}

	var (
		cfg             = ds.GetConfig()
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
		repository = shorten.NewRepository(dbPool, logger)
		service    = NewService(repository)
	)

	t.Run("Create new short url with GenerateShortenUrl()", func(t *testing.T) {
		var (
			inputShorten = ds.InputShorten{
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
			inputShorten = ds.InputShorten{
				ShortenUrl: mo.Some("google"),
				OriginUrl:  "https://google.com/",
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

		count, errGen := service.DeleteShorten(context.Background(), generateShorten)
		assert.NoError(t, errGen)
		assert.NotZero(t, count)

		count, errCustom := service.DeleteShorten(context.Background(), customShorten)
		assert.NoError(t, errCustom)
		assert.NotZero(t, count)
	})

	// delete dir with logs
	err = os.RemoveAll(filepath.Join(testDir, "logs"))
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

}
