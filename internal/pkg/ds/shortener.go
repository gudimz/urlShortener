package ds

import (
	"time"

	"github.com/samber/mo"
)

type Shorten struct {
	ShortURL    string    `json:"short_url"`
	OriginURL   string    `json:"origin_url"`
	Visits      int64     `json:"visits"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

type InputShorten struct {
	ShortenURL mo.Option[string]
	OriginURL  string
}
