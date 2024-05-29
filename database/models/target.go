package models

type TargetVisibility string

const (
	TargetVisibilityPublic  TargetVisibility = "PUBLIC"
	TargetVisibilityPrivate TargetVisibility = "PRIVATE"
	TargetVisibilityUnknown TargetVisibility = "UNKNOWN"
)

type TargetState string

const (
	TargetStateActive  TargetState = "ACTIVE"
	TargetStatePaused  TargetState = "PAUSED"
	TargetStateUnknown TargetState = "UNKNOWN"
)

type Target struct {
	CreatedAt *Time `db:"created_at"`
	UpdatedAt *Time `db:"updated_at"`

	Id         string           `db:"id"`
	Name       string           `db:"name"`
	Group      string           `db:"group"`
	Visibility TargetVisibility `db:"visibility"`
	State      TargetState      `db:"state"`
	Metadata   string           `db:"metadata"`
}

type TargetStatus string

const (
	TargetStatusSuccess TargetStatus = "SUCCESS"
	TargetStatusFailure TargetStatus = "FAILURE"
	TargetStatusUnknown TargetStatus = "UNKNOWN"
)

type TargetHistory struct {
	CreatedAt *Time `db:"created_at"`

	TargetId      string `db:"target_id"`
	WorkerGroupId string `db:"worker_group_id"`
	CheckId       string `db:"check_id"`

	Status TargetStatus `db:"status"`
	Note   string       `db:"note"`
}
