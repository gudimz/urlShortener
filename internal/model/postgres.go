package model

import "time"

type dbShorten struct {
	ID            string    `bson:"id"`
	OriginUrl     string    `bson:"origin_url"`
	VisitsCounter int64     `bson:"visits_counter"`
	DateCreated   time.Time `bson:"date_created"`
	DateUpdated   time.Time `bson:"date_updated"`
}
