package models

import (
	"time"
)

type OAuth2State struct {
	State     string    `db:"state"`
	ExpiresAt time.Time `db:"expires_at"`
}

const (
	MonitorSuccess string = "SUCCESS"
	MonitorFailure string = "FAILURE"
	MonitorError   string = "ERROR"
	MonitorUnknown string = "UNKNOWN"
)

type Monitor struct {
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	Slug string `db:"slug"`
	Name string `db:"name"`

	Schedule string `db:"schedule"`
	Script   string `db:"script"`
}

type MonitorWithWorkerGroups struct {
	Monitor

	// List of worker group names
	WorkerGroups []string
}

type MonitorHistory struct {
	CreatedAt time.Time `db:"created_at"`

	MonitorSlug string `db:"monitor_slug"`
	Status      string `db:"status"`
	Note        string `db:"note"`
}

type WorkerGroup struct {
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	Slug string `db:"slug"`
	Name string `db:"name"`
}

type WorkerGroupWithMonitors struct {
	WorkerGroup

	// List of worker group names
	Monitors []string
}
