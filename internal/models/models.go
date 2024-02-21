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
	Slug string `gorm:"unique"`
	Name string `gorm:"unique" validate:"required"`

	Schedule     string         `validate:"required,cron"`
	WorkerGroups pq.StringArray `gorm:"type:text[]"`

	Script string `validate:"required"`

	History []HealthcheckHistory `gorm:"foreignKey:ID"`
}

type Cronjob struct {
	gorm.Model
	Slug     string `gorm:"unique"`
	Name     string `gorm:"unique"`
	Schedule string
	Buffer   int
}

type HealthcheckHistory struct {
	gorm.Model
	Healthcheck Healthcheck `gorm:"foreignkey:ID"`
	Status      string
}

type CronjobHistory struct {
	gorm.Model
	Cronjob Cronjob `gorm:"foreignkey:ID"`
	Status  string
}
