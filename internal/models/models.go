package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuth2State struct {
	State  string `gorm:"primary_key"`
	Expiry time.Time
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
	Name string `gorm:"unique"`

	Schedule     string
	WorkerGroups []WorkerGroup `gorm:"many2many:monitor_worker_groups;"`

	Script string `validate:"required"`

	History []MonitorHistory `gorm:"foreignKey:Monitor"`
}

type MonitorHistory struct {
	gorm.Model
	Monitor uint
	Status  string
	Note    string
}

type WorkerGroup struct {
	gorm.Model
	Name string `gorm:"unique"`
	Slug string `gorm:"unique"`

	Monitors []Monitor `gorm:"many2many:monitor_worker_groups;"`
}
