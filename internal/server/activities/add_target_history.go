package activities

import (
	"context"

	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/internal/temporal"
)

func (a *Activities) AddTargetHistory(ctx context.Context, param temporal.ActivityAddTargetHistoryParam) (*temporal.ActivityAddTargetHistoryResult, error) {

	status := models.TargetStatusUnknown
	if param.Status == temporal.AddTargetHistoryStatusSuccess {
		status = models.TargetStatusSuccess
	}
	if param.Status == temporal.AddTargetHistoryStatusFailure {
		status = models.TargetStatusFailure
	}

	err := services.AddHistoryForTarget(ctx, a.db, &models.TargetHistory{
		TargetId: param.Target.Id,
		Status:   status,
		Note:     param.Note,
	})

	return &temporal.ActivityAddTargetHistoryResult{}, err
}
