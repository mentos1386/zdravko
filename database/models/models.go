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

type TriggerState string

const (
	TriggerStateActive  TriggerState = "ACTIVE"
	TriggerStatePaused  TriggerState = "PAUSED"
	TriggerStateUnknown TriggerState = "UNKNOWN"
)

type Trigger struct {
	CreatedAt *Time `db:"created_at"`
	UpdatedAt *Time `db:"updated_at"`

	Id     string `db:"id"`
	Name   string `db:"name"`
	Script string `db:"script"`
}
