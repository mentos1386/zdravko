package models

type CheckState string

const (
	CheckStateActive  CheckState = "ACTIVE"
	CheckStatePaused  CheckState = "PAUSED"
	CheckStateUnknown CheckState = "UNKNOWN"
)

type Check struct {
	CreatedAt *Time `db:"created_at"`
	UpdatedAt *Time `db:"updated_at"`

	Id   string `db:"id"`
	Name string `db:"name"`

	Schedule string `db:"schedule"`
	Script   string `db:"script"`
	Filter   string `db:"filter"`
}

type CheckWithWorkerGroups struct {
	Check
	// List of worker group names
	WorkerGroups []string
}
