package models

import "time"

type OAuth2State struct {
	State  string `gorm:"primary_key"`
	Expiry time.Time
}

type Healthcheck struct {
	ID               uint `gorm:"primary_key"`
	Name             string
	Status           string // UP, DOWN
	UptimePercentage float64
}
