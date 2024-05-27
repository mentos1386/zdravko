package activities

import (
	"context"

	"github.com/dop251/goja"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/internal/temporal"
	"github.com/mentos1386/zdravko/pkg/script"
	"gopkg.in/yaml.v3"
)

func (a *Activities) TargetsFilter(ctx context.Context, param temporal.ActivityTargetsFilterParam) (*temporal.ActivityTargetsFilterResult, error) {
	a.logger.Info("TargetsFilter", "filter", param.Filter)

	allTargets, err := services.GetTargets(ctx, a.db)
	if err != nil {
		return nil, err
	}
	filteredTargets := make([]*temporal.Target, 0)

	program, err := goja.Compile("filter", script.UnescapeString(param.Filter), false)
	if err != nil {
		return nil, err
	}

	for _, target := range allTargets {
		if target.State == models.TargetStatePaused {
			continue
		}

		vm := goja.New()
		vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

		var metadata map[string]interface{}
		err := yaml.Unmarshal([]byte(target.Metadata), &metadata)
		if err != nil {
			return nil, err
		}

		a.logger.Info("TargetsFilter", "target", target)

		targetWithMedatada := &struct {
			Id       string
			Name     string
			Group    string
			Metadata map[string]interface{}
		}{
			Id:       target.Id,
			Name:     target.Name,
			Group:    target.Group,
			Metadata: metadata,
		}

		err = vm.Set("target", targetWithMedatada)
		if err != nil {
			return nil, err
		}

		value, err := vm.RunProgram(program)
		if err != nil {
			return nil, err
		}
		if value.Export().(bool) {
			filteredTargets = append(filteredTargets, &temporal.Target{
				Id:       target.Id,
				Name:     target.Name,
				Group:    target.Group,
				Metadata: target.Metadata,
			})
		}
	}

	return &temporal.ActivityTargetsFilterResult{
		Targets: filteredTargets,
	}, nil
}
