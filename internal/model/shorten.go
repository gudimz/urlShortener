package model

import "time"

type shorten struct {
	ID            string    `json:"id"`
	OriginUrl     string    `json:"origin_url"`
	VisitsCounter int64     `json:"visits_counter"`
	DateCreated   time.Time `json:"date_created"`
	DateUpdated   time.Time `json:"date_updated"`
}
