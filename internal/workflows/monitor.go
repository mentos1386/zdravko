package workflows

import (
	"sort"
	"time"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"go.temporal.io/sdk/workflow"
)

type MonitorWorkflowParam struct {
	Script       string
	Slug         string
	WorkerGroups []string
}

func (w *Workflows) MonitorWorkflowDefinition(ctx workflow.Context, param MonitorWorkflowParam) error {
	workerGroups := param.WorkerGroups
	sort.Strings(workerGroups)

	for _, workerGroup := range workerGroups {
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 60 * time.Second,
			TaskQueue:           workerGroup,
		})

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
			Slug:        param.Slug,
			Status:      status,
			Note:        monitorResult.Note,
			WorkerGroup: workerGroup,
		}

		var historyResult *activities.MonitorAddToHistoryResult
		err = workflow.ExecuteActivity(ctx, w.activities.MonitorAddToHistory, historyParam).Get(ctx, &historyResult)
		if err != nil {
			return err
		}
	}

	return nil
}
