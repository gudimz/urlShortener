package model

import "time"

type Shorten struct {
	ID          string    `json:"id"`
	OriginUrl   string    `json:"origin_url"`
	Visits      int64     `json:"visits"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}
