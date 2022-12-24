package model

import "time"

type shorten struct {
	ID            string    `json:"id"`
	OriginUrl     string    `json:"origin_url"`
	VisitsCounter int64     `json:"visits_counter"`
	Created       time.Time `json:"created"` //Date of creation
	Updated       time.Time `json:"updated"` //Date of updated
}
