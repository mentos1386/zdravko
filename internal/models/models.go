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

const (
	MonitorSuccess string = "SUCCESS"
	MonitorFailure string = "FAILURE"
	MonitorError   string = "ERROR"
	MonitorUnknown string = "UNKNOWN"
)

type Monitor struct {
	gorm.Model
	Slug string `gorm:"unique"`
	Name string `gorm:"unique" validate:"required"`

	Schedule     string         `validate:"required,cron"`
	WorkerGroups pq.StringArray `gorm:"type:text[]"`

	Script string `validate:"required"`

	History []MonitorHistory `gorm:"foreignKey:Monitor"`
}

type Cronjob struct {
	gorm.Model
	Slug     string `gorm:"unique"`
	Name     string `gorm:"unique"`
	Schedule string
	Buffer   int
}

type MonitorHistory struct {
	gorm.Model
	Monitor uint
	Status      string
	Note        string
}

type CronjobHistory struct {
	gorm.Model
	Cronjob Cronjob `gorm:"foreignkey:ID"`
	Status  string
}
