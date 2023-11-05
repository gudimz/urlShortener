package models

import (
	"time"

	"github.com/gudimz/urlShortener/internal/pkg/ds"
)

type DbShorten struct {
	ShortUrl    string    `bson:"short_url"`
	OriginUrl   string    `bson:"origin_url"`
	Visits      int64     `bson:"visits"`
	DateCreated time.Time `bson:"date_created"`
	DateUpdated time.Time `bson:"date_updated"`
}

func ModelFromDbShorten(dbShorten *DbShorten) *ds.Shorten {
	return &ds.Shorten{
		ShortUrl:    dbShorten.ShortUrl,
		OriginUrl:   dbShorten.OriginUrl,
		Visits:      dbShorten.Visits,
		DateCreated: dbShorten.DateCreated,
		DateUpdated: dbShorten.DateUpdated,
	}
}

func DbShortenFromModel(shorten *ds.Shorten) *DbShorten {
	return &DbShorten{
		ShortUrl:  shorten.ShortUrl,
		OriginUrl: shorten.OriginUrl,
	}
}
