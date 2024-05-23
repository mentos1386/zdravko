package activities

import (
	"context"

	"github.com/hashicorp/mql"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/internal/temporal"
)

func (a *Activities) TargetsFilter(ctx context.Context, param temporal.ActivityTargetsFilterParam) (*temporal.ActivityTargetsFilterResult, error) {
	a.logger.Info("TargetsFilter", "filter", param.Filter)
	f, err := mql.Parse(param.Filter)
	if err != nil {
		return nil, err
	}
	a.logger.Info("TargetsFilter", "filter", f)

	// TODO: Parse filter.
	targets, err := services.GetTargetsWithFilter(ctx, a.db, param.Filter)
	if err != nil {
		return nil, err
	}

	return &temporal.ActivityTargetsFilterResult{
		Targets: targets,
	}, nil
}
