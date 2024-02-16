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
	Slug             string `gorm:"unique"`
	Name             string `gorm:"unique"`
	Status           string // UP, DOWN
	UptimePercentage float64
	Schedule         string
}

type HealthcheckHttp struct {
	gorm.Model
	Healthcheck
	Url    string
	Method string
}

type HealthcheckTcp struct {
	gorm.Model
	Healthcheck
	Hostname string
	Port     int
}

type Cronjob struct {
	gorm.Model
	Slug     string `gorm:"unique"`
	Name     string `gorm:"unique"`
	Schedule string
	Buffer   int
}

type HealthcheckHttpHistory struct {
	gorm.Model
	HealthcheckHTTP HealthcheckHttp `gorm:"foreignkey:ID"`
	Status          string
}

type HealthcheckTcpHistory struct {
	gorm.Model
	HealthcheckTCP HealthcheckTcp `gorm:"foreignkey:ID"`
	Status         string
}

type CronjobHistory struct {
	gorm.Model
	Cronjob Cronjob `gorm:"foreignkey:ID"`
	Status  string
}
