package activities

import (
	"context"

	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
)

type TargetsFilterParam struct {
	Filter string
}

type TargetsFilterResult struct {
	Targets []*models.Target
}

const TargetsFilterName = "TARGETS_FILTER"

func (a *Activities) TargetsFilter(ctx context.Context, param TargetsFilterParam) (*TargetsFilterResult, error) {
	a.logger.Info("TargetsFilter", "filter", param.Filter)
	// TODO: Parse filter.
	targets, err := services.GetTargetsWithFilter(ctx, a.db, param.Filter)
	if err != nil {
		return nil, err
	}

	return &TargetsFilterResult{
		Targets: targets,
	}, nil
}
