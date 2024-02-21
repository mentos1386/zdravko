package workflows

import (
	"code.tjo.space/mentos1386/zdravko/internal/activities"
)

type Workflows struct {
	activities *activities.Activities
}

func NewWorkflows(a *activities.Activities) *Workflows {
	return &Workflows{activities: a}
}
