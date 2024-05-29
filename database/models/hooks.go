package models

type HookState string

const (
	HookStateActive  HookState = "ACTIVE"
	HookStatePaused  HookState = "PAUSED"
	HookStateUnknown HookState = "UNKNOWN"
)

type Hook struct {
	CreatedAt *Time `db:"created_at"`
	UpdatedAt *Time `db:"updated_at"`

	Id   string `db:"id"`
	Name string `db:"name"`

	Schedule string `db:"schedule"`
	Script   string `db:"script"`
}
