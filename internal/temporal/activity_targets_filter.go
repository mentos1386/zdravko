package temporal

import (
	"github.com/mentos1386/zdravko/database/models"
)

type ActivityTargetsFilterParam struct {
	Filter string
}

type ActivityTargetsFilterResult struct {
	Targets []*models.Target
}

const ActivityTargetsFilterName = "TARGETS_FILTER"
