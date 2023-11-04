package shorten

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gudimz/urlShortener/internal/app/repository/psql/models"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/gudimz/urlShortener/pkg/postgres"
)

const (
	queryToCreate = `
						INSERT INTO shorten_urls
							(short_url, origin_url, visits, date_created, date_updated)
						VALUES
							($1, $2, $3, $4, $5)
	`
	queryToGetAll = `
						SELECT
							short_url, origin_url, visits, date_created, date_updated
						FROM
							shorten_urls
	`
	queryToGetByShortUrl    = queryToGetAll + ` WHERE short_url = $1`
	queryToDeleteByShortUrl = `DELETE FROM shorten_urls WHERE short_url = $1`
	queryToUpdateByShortUrl = `
						UPDATE shorten_urls
							SET visits = visits + 1, date_updated = $1
						WHERE short_url = $2
	`
)

type Repository struct {
	db     postgres.Client
	logger *logging.Logger
}

func NewRepository(db postgres.Client, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func queryForLogger(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *Repository) CreateShorten(ctx context.Context, ms ds.Shorten) error {
	q := queryToCreate
	dbm := models.DbShortenFromModel(ms)
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	_, err := r.db.Exec(
		ctx,
		q,
		dbm.ShortUrl,
		dbm.OriginUrl,
		dbm.Visits,
		dbm.DateCreated,
		dbm.DateUpdated)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}

func (r *Repository) GetShorten(ctx context.Context, shortUrl string) (*ds.Shorten, error) {
	q := queryToGetByShortUrl
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	var dbShorten models.DbShorten
	err := r.db.QueryRow(ctx, q, shortUrl).Scan(
		&dbShorten.ShortUrl,
		&dbShorten.OriginUrl,
		&dbShorten.Visits,
		&dbShorten.DateCreated,
		&dbShorten.DateUpdated)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return models.ModelFromDbShorten(dbShorten), nil
}

func (r *Repository) DeleteShorten(ctx context.Context, shortUrl string) (int64, error) {
	q := queryToDeleteByShortUrl
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	res, err := r.db.Exec(ctx, q, shortUrl)
	if err != nil {
		r.logger.Error(err)
		return 0, err
	}
	count := res.RowsAffected()
	return count, nil
}

func (r *Repository) UpdateShorten(ctx context.Context, shortUrl string) error {
	q := queryToUpdateByShortUrl
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	_, err := r.db.Exec(ctx, q, time.Now().UTC(), shortUrl)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}
