package activities

import (
	"context"

	"github.com/dop251/goja"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/internal/temporal"
	"github.com/mentos1386/zdravko/pkg/script"
)

func (a *Activities) TargetsFilter(ctx context.Context, param temporal.ActivityTargetsFilterParam) (*temporal.ActivityTargetsFilterResult, error) {
	a.logger.Info("TargetsFilter", "filter", param.Filter)

	allTargets, err := services.GetTargets(ctx, a.db)
	if err != nil {
		return nil, err
	}
	filteredTargets := make([]*models.Target, 0)

	for _, target := range allTargets {
		vm := goja.New()

		err = vm.Set("target", target)
		if err != nil {
			return nil, err
		}

		value, err := vm.RunString(script.UnescapeString(param.Filter))
		if err != nil {
			return nil, err
		}
		if value.Export().(bool) {
			filteredTargets = append(filteredTargets, target)
		}
	}

	return &temporal.ActivityTargetsFilterResult{
		Targets: filteredTargets,
	}, nil
}
