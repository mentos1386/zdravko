package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuth2State struct {
	State  string `gorm:"primary_key"`
	Expiry time.Time
}

type Healthcheck struct {
	gorm.Model
	Name             string `gorm:"unique"`
	Status           string // UP, DOWN
	UptimePercentage float64
	Schedule         string
}

type HealthcheckHTTP struct {
	gorm.Model
	Healthcheck
	URL    string
	Method string
}

type HealthcheckTCP struct {
	gorm.Model
	Healthcheck
	Hostname string
	Port     int
}

type Cronjob struct {
	gorm.Model
	Name     string `gorm:"unique"`
	Schedule string
	Buffer   int
}
