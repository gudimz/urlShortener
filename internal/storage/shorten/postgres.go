package shorten

import (
	"context"
	"fmt"
	"github.com/gudimz/urlShortener/internal/db/postgres"
	"github.com/gudimz/urlShortener/internal/model"
	"github.com/gudimz/urlShortener/pkg/logging"
	"strings"
	"time"
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
						WHERE id = $2
	`
)

type storage struct {
	db     postgres.Client
	logger *logging.Logger
}

func NewStorage(db postgres.Client, logger *logging.Logger) *storage {
	return &storage{
		db:     db,
		logger: logger,
	}
}

func queryForLogger(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (p *storage) CreateShorten(ctx context.Context, ms model.Shorten) error {
	q := queryToCreate
	dbm := dbShortenFromModel(ms)
	p.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	_, err := p.db.Exec(
		ctx,
		q,
		dbm.ShortUrl,
		dbm.OriginUrl,
		dbm.Visits,
		dbm.DateCreated,
		dbm.DateUpdated)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	return nil
}

func (p *storage) GetShorten(ctx context.Context, shortUrl string) (*model.Shorten, error) {
	q := queryToGetByShortUrl
	p.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	var dbShorten model.DbShorten
	err := p.db.QueryRow(ctx, q, shortUrl).Scan(
		&dbShorten.ShortUrl,
		&dbShorten.OriginUrl,
		&dbShorten.Visits,
		&dbShorten.DateCreated,
		&dbShorten.DateUpdated)
	if err != nil {
		p.logger.Error(err)
		return nil, err
	}
	return modelFromDbShorten(dbShorten), nil
}

func (p *storage) DeleteShorten(ctx context.Context, shortUrl string) error {
	q := queryToDeleteByShortUrl
	p.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	_, err := p.db.Exec(ctx, q, shortUrl)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	return nil
}

func (p *storage) UpdateShorten(ctx context.Context, shortUrl string) error {
	q := queryToUpdateByShortUrl
	p.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	_, err := p.db.Exec(ctx, q, time.Now().UTC(), shortUrl)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	return nil
}

func dbShortenFromModel(shorten model.Shorten) model.DbShorten {
	return model.DbShorten{
		ShortUrl:    shorten.ShortUrl,
		OriginUrl:   shorten.OriginUrl,
		Visits:      shorten.Visits,
		DateCreated: shorten.DateCreated,
		DateUpdated: shorten.DateUpdated,
	}
}

func modelFromDbShorten(dbShorten model.DbShorten) *model.Shorten {
	return &model.Shorten{
		ShortUrl:    dbShorten.ShortUrl,
		OriginUrl:   dbShorten.OriginUrl,
		Visits:      dbShorten.Visits,
		DateCreated: dbShorten.DateCreated,
		DateUpdated: dbShorten.DateUpdated,
	}
}
