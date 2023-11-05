package models

import (
	"time"

	"github.com/gudimz/urlShortener/internal/pkg/ds"
)

type DBShorten struct {
	ShortURL    string    `bson:"short_url"`
	OriginURL   string    `bson:"origin_url"`
	Visits      int64     `bson:"visits"`
	DateCreated time.Time `bson:"date_created"`
	DateUpdated time.Time `bson:"date_updated"`
}

func ModelFromDBShorten(dbShorten *DBShorten) *ds.Shorten {
	return &ds.Shorten{
		ShortURL:    dbShorten.ShortURL,
		OriginURL:   dbShorten.OriginURL,
		Visits:      dbShorten.Visits,
		DateCreated: dbShorten.DateCreated,
		DateUpdated: dbShorten.DateUpdated,
	}
}

func DBShortenFromModel(shorten *ds.Shorten) *DBShorten {
	return &DBShorten{
		ShortURL:  shorten.ShortURL,
		OriginURL: shorten.OriginURL,
	}
}
