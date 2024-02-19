package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type OAuth2State struct {
	State  string `gorm:"primary_key"`
	Expiry time.Time
}

type Worker struct {
	gorm.Model
	Name   string `gorm:"unique" validate:"required"`
	Slug   string `gorm:"unique"`
	Group  string `validate:"required"`
	Status string
}

type Healthcheck struct {
	gorm.Model
	Slug             string `gorm:"unique"`
	Name             string `gorm:"unique" validate:"required"`
	Status           string // UP, DOWN
	UptimePercentage float64
	Schedule         string         `validate:"required,cron"`
	WorkerGroups     pq.StringArray `gorm:"type:text[]"`
}

type HealthcheckHttp struct {
	gorm.Model
	Healthcheck
	Url    string `validate:"required,url"`
	Method string `validate:"required,oneof=GET POST"`
}

type HealthcheckTcp struct {
	gorm.Model
	Healthcheck
	Hostname string `validate:"required,hostname"`
	Port     int    `validate:"required,gte=1,lte=65535"`
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
