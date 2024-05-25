package activities

import (
	"context"
	"log/slog"

	"github.com/mentos1386/zdravko/internal/temporal"
	"github.com/mentos1386/zdravko/pkg/k6"
	"github.com/mentos1386/zdravko/pkg/k6/zdravko"
	"github.com/mentos1386/zdravko/pkg/script"
	"gopkg.in/yaml.v3"
)

func (a *Activities) Check(ctx context.Context, param temporal.ActivityCheckParam) (*temporal.ActivityCheckResult, error) {
	execution := k6.NewExecution(slog.Default(), script.UnescapeString(param.Script))

	var metadata map[string]interface{}
	err := yaml.Unmarshal([]byte(param.Target.Metadata), &metadata)
	if err != nil {
		return nil, err
	}

	ctx = zdravko.WithZdravkoContext(ctx, zdravko.Context{
		Target: zdravko.Target{
			Name:     param.Target.Name,
			Group:    param.Target.Group,
			Metadata: metadata,
		},
	})

	result, err := execution.Run(ctx)
	if err != nil {
		return nil, err
	}

	return &temporal.ActivityCheckResult{Success: result.Success, Note: result.Note}, nil
}
