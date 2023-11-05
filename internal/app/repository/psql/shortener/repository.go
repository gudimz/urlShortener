package shortener

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gudimz/urlShortener/internal/app/repository/psql/models"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
	"github.com/gudimz/urlShortener/pkg/logger"
	"github.com/gudimz/urlShortener/pkg/postgres"
)

type Repository struct {
	db  postgres.Client
	log *logger.Log
}

func NewRepository(db postgres.Client, log *logger.Log) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func queryForLogger(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *Repository) CreateShorten(ctx context.Context, ms *ds.Shorten) (*models.DbShorten, error) {
	op := "CreateShorten"
	q := `INSERT INTO shorten_urls
			(short_url, origin_url)
		VALUES
			($1, $2)
		RETURNING *;`

	dbm := models.DbShortenFromModel(ms)
	r.log.Debug(op, zap.String("query", queryForLogger(q)))

	var newShorten models.DbShorten
	err := r.db.QueryRow(
		ctx,
		q,
		dbm.ShortUrl,
		dbm.OriginUrl,
	).Scan(&newShorten.ShortUrl, &newShorten.OriginUrl, &newShorten.Visits, &newShorten.DateCreated, &newShorten.DateUpdated)

	if err != nil {
		r.log.Error(op, zap.Error(err))
		return nil, err
	}

	return &newShorten, nil
}

func (r *Repository) GetShorten(ctx context.Context, shortUrl string) (*models.DbShorten, error) {
	op := "GetShorten"
	q := `SELECT
			short_url, origin_url, visits, date_created, date_updated
		FROM
			shorten_urls
		WHERE short_url = $1;`
	r.log.Debug(op, zap.String("query", queryForLogger(q)))

	var dbShorten models.DbShorten
	err := r.db.QueryRow(ctx, q, shortUrl).Scan(
		&dbShorten.ShortUrl,
		&dbShorten.OriginUrl,
		&dbShorten.Visits,
		&dbShorten.DateCreated,
		&dbShorten.DateUpdated)
	if err != nil {
		r.log.Error(op, zap.Error(err))
		return nil, err
	}

	return &dbShorten, nil
}

func (r *Repository) DeleteShorten(ctx context.Context, shortUrl string) (int64, error) {
	op := "DeleteShorten"
	q := `DELETE FROM
			shorten_urls
		WHERE short_url = $1;`
	r.log.Debug(op, zap.String("query", queryForLogger(q)))
	res, err := r.db.Exec(ctx, q, shortUrl)
	if err != nil {
		r.log.Error(op, zap.Error(err))
		return 0, err
	}
	count := res.RowsAffected()
	return count, nil
}

func (r *Repository) UpdateShorten(ctx context.Context, shortUrl string) error {
	op := "DeleteShorten"
	q := `UPDATE 
			shorten_urls
		SET
			visits = visits + 1, date_updated = $1
		WHERE
			short_url = $2
		RETURNING *;`

	r.log.Debug(op, zap.String("query", queryForLogger(q)))
	_, err := r.db.Exec(ctx, q, time.Now().UTC(), shortUrl)
	if err != nil {
		r.log.Error(op, zap.Error(err))
		return err
	}

	return nil
}
