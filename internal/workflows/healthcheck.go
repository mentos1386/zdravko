package workflows

import (
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"go.temporal.io/sdk/workflow"
)

type HealthcheckWorkflowParam struct {
	Script string
	Slug   string
}

func (w *Workflows) HealthcheckWorkflowDefinition(ctx workflow.Context, param HealthcheckWorkflowParam) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	heatlcheckParam := activities.HealtcheckParam{
		Script: param.Script,
	}

	var healthcheckResult *activities.HealthcheckResult
	err := workflow.ExecuteActivity(ctx, w.activities.Healthcheck, heatlcheckParam).Get(ctx, &healthcheckResult)
	if err != nil {
		return err
	}

	status := models.HealthcheckFailure
	if healthcheckResult.Success {
		status = models.HealthcheckSuccess
	}

	historyParam := activities.HealtcheckAddToHistoryParam{
		Slug:   param.Slug,
		Status: status,
		Note:   healthcheckResult.Note,
	}

	var historyResult *activities.HealthcheckAddToHistoryResult
	err = workflow.ExecuteActivity(ctx, w.activities.HealthcheckAddToHistory, historyParam).Get(ctx, &historyResult)
	if err != nil {
		return err
	}

	return nil
}
