package shorten

import (
	"context"
	"fmt"
	"github.com/gudimz/urlShortener/internal/db/postgres"
	"github.com/gudimz/urlShortener/internal/model"
	"github.com/gudimz/urlShortener/internal/shorten"
	"github.com/gudimz/urlShortener/pkg/logging"
	"strings"
	"time"
)

const (
	queryToCreate = `
						INSERT INTO shorten_urls
							(id, origin_url, visits, date_created, date_updated)
						VALUES
							($1, $2, $3, $4, $5)
	`
	queryToGetAll = `
						SELECT
							id, origin_url, visits, date_created, date_updated
						FROM
							shorten_urls
	`
	queryToGetByID    = queryToGetAll + ` WHERE id = $1`
	queryToDeleteByID = `DELETE FROM shorten_urls WHERE id = $1`
	queryToUpdateByID = `
						UPDATE shorten_urls
							SET visits = visits + 1, date_updated = $1
						WHERE id = $2
	`
)

type storage struct {
	db     postgres.Client
	logger *logging.Logger
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
		dbm.ID,
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

func (p *storage) GetShorten(ctx context.Context, id string) (*model.Shorten, error) {
	q := queryToGetByID
	p.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	var dbShorten model.DbShorten
	err := p.db.QueryRow(ctx, q, id).Scan(
		&dbShorten.ID,
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

func (p *storage) DeleteShorten(ctx context.Context, id string) error {
	q := queryToDeleteByID
	p.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	_, err := p.db.Exec(ctx, q, id)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	return nil
}

func (p *storage) UpdateShorten(ctx context.Context, id string) error {
	q := queryToUpdateByID
	p.logger.Trace(fmt.Sprintf("SQL Query: %s", queryForLogger(q)))
	_, err := p.db.Exec(ctx, q, time.Now().UTC(), id)
	if err != nil {
		p.logger.Error(err)
		return err
	}
	return nil
}

func dbShortenFromModel(shorten model.Shorten) model.DbShorten {
	return model.DbShorten{
		ID:          shorten.ID,
		OriginUrl:   shorten.OriginUrl,
		Visits:      shorten.Visits,
		DateCreated: shorten.DateCreated,
		DateUpdated: shorten.DateUpdated,
	}
}

func modelFromDbShorten(dbShorten model.DbShorten) *model.Shorten {
	return &model.Shorten{
		ID:          dbShorten.ID,
		OriginUrl:   dbShorten.OriginUrl,
		Visits:      dbShorten.Visits,
		DateCreated: dbShorten.DateCreated,
		DateUpdated: dbShorten.DateUpdated,
	}
}

func NewStorage(db postgres.Client, logger *logging.Logger) shorten.Storage {
	return &storage{
		db:     db,
		logger: logger,
	}
}
