package model

import "time"

type DbShorten struct {
	ID          string    `bson:"id"`
	OriginUrl   string    `bson:"origin_url"`
	Visits      int64     `bson:"visits"`
	DateCreated time.Time `bson:"date_created"`
	DateUpdated time.Time `bson:"date_updated"`
}
