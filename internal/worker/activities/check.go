package activities

import (
	"context"
	"log/slog"

	"github.com/mentos1386/zdravko/internal/temporal"
	"github.com/mentos1386/zdravko/pkg/k6"
	_ "github.com/mentos1386/zdravko/pkg/k6/zdravko"
	"github.com/mentos1386/zdravko/pkg/script"
)

func (a *Activities) Check(ctx context.Context, param temporal.ActivityCheckParam) (*temporal.ActivityCheckResult, error) {
	execution := k6.NewExecution(slog.Default(), script.UnescapeString(param.Script))

	result, err := execution.Run(ctx)
	if err != nil {
		return nil, err
	}

	return &temporal.ActivityCheckResult{Success: result.Success, Note: result.Note}, nil
}
