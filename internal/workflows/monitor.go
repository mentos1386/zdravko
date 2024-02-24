package workflows

import (
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"go.temporal.io/sdk/workflow"
)

type MonitorWorkflowParam struct {
	Script string
	Slug   string
}

func (w *Workflows) MonitorWorkflowDefinition(ctx workflow.Context, param MonitorWorkflowParam) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	heatlcheckParam := activities.HealtcheckParam{
		Script: param.Script,
	}

	var monitorResult *activities.MonitorResult
	err := workflow.ExecuteActivity(ctx, w.activities.Monitor, heatlcheckParam).Get(ctx, &monitorResult)
	if err != nil {
		return err
	}

	status := models.MonitorFailure
	if monitorResult.Success {
		status = models.MonitorSuccess
	}

	historyParam := activities.HealtcheckAddToHistoryParam{
		Slug:   param.Slug,
		Status: status,
		Note:   monitorResult.Note,
	}

	var historyResult *activities.MonitorAddToHistoryResult
	err = workflow.ExecuteActivity(ctx, w.activities.MonitorAddToHistory, historyParam).Get(ctx, &historyResult)
	if err != nil {
		return err
	}

	return nil
}
