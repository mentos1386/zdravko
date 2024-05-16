package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	Time time.Time
}

// rfc3339Milli is like time.RFC3339Nano, but with millisecond precision, and fractional seconds do not have trailing
// zeros removed.
const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

// Value satisfies driver.Valuer interface.
func (t *Time) Value() (driver.Value, error) {
	return t.Time.UTC().Format(rfc3339Milli), nil
}

// Scan satisfies sql.Scanner interface.
func (t *Time) Scan(src any) error {
	if src == nil {
		return nil
	}

	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("error scanning time, got %+v", src)
	}

	parsedT, err := time.Parse(rfc3339Milli, s)
	if err != nil {
		return err
	}

	t.Time = parsedT.UTC()

	return nil
}

type OAuth2State struct {
	State     string `db:"state"`
	ExpiresAt *Time  `db:"expires_at"`
}

type CheckStatus string

const (
	CheckSuccess CheckStatus = "SUCCESS"
	CheckFailure CheckStatus = "FAILURE"
	CheckError   CheckStatus = "ERROR"
	CheckUnknown CheckStatus = "UNKNOWN"
)

type Check struct {
	CreatedAt *Time `db:"created_at"`
	UpdatedAt *Time `db:"updated_at"`

	Id    string `db:"id"`
	Name  string `db:"name"`
	Group string `db:"group"`

	Schedule string `db:"schedule"`
	Script   string `db:"script"`
}

type CheckWithWorkerGroups struct {
	Check

	// List of worker group names
	WorkerGroups []string
}

type CheckHistory struct {
	CreatedAt *Time `db:"created_at"`

	CheckId string        `db:"check_id"`
	Status    CheckStatus `db:"status"`
	Note      string        `db:"note"`

	WorkerGroupId   string `db:"worker_group_id"`
	WorkerGroupName string `db:"worker_group_name"`
}

type WorkerGroup struct {
	CreatedAt *Time `db:"created_at"`
	UpdatedAt *Time `db:"updated_at"`

	Id   string `db:"id"`
	Name string `db:"name"`
}

type WorkerGroupWithChecks struct {
	WorkerGroup

	// List of worker group names
	Checks []string
}

type TriggerStatus string

const (
	TriggerSuccess TriggerStatus = "SUCCESS"
	TriggerFailure TriggerStatus = "FAILURE"
	TriggerError   TriggerStatus = "ERROR"
	TriggerUnknown TriggerStatus = "UNKNOWN"
)

type Trigger struct {
	CreatedAt *Time `db:"created_at"`
	UpdatedAt *Time `db:"updated_at"`

	Id     string        `db:"id"`
	Name   string        `db:"name"`
	Script string        `db:"script"`
	Status TriggerStatus `db:"status"`
}

type TriggerHistory struct {
	CreatedAt *Time `db:"created_at"`

	TriggerId string        `db:"trigger_id"`
	Status    TriggerStatus `db:"status"`
	Note      string        `db:"note"`
}
