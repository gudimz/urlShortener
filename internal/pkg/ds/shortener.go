package ds

import (
	"time"

	"github.com/samber/mo"
)

type Shorten struct {
	ShortUrl    string    `json:"short_url"`
	OriginUrl   string    `json:"origin_url"`
	Visits      int64     `json:"visits"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

type InputShorten struct {
	ShortenUrl mo.Option[string]
	OriginUrl  string
}
