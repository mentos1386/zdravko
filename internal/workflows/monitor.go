package workflows

import (
	"sort"
	"time"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"go.temporal.io/sdk/workflow"
)

type MonitorWorkflowParam struct {
	Script         string
	MonitorId      string
	WorkerGroupIds []string
}

func (w *Workflows) MonitorWorkflowDefinition(ctx workflow.Context, param MonitorWorkflowParam) (models.MonitorStatus, error) {
	workerGroupIds := param.WorkerGroupIds
	sort.Strings(workerGroupIds)

	for _, workerGroupId := range workerGroupIds {
		ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 60 * time.Second,
			TaskQueue:           workerGroupId,
		})

		heatlcheckParam := activities.HealtcheckParam{
			Script: param.Script,
		}

		var monitorResult *activities.MonitorResult
		err := workflow.ExecuteActivity(ctx, w.activities.Monitor, heatlcheckParam).Get(ctx, &monitorResult)
		if err != nil {
			return models.MonitorUnknown, err
		}

		status := models.MonitorFailure
		if monitorResult.Success {
			status = models.MonitorSuccess
		}

		historyParam := activities.HealtcheckAddToHistoryParam{
			MonitorId:     param.MonitorId,
			Status:        status,
			Note:          monitorResult.Note,
			WorkerGroupId: workerGroupId,
		}

		var historyResult *activities.MonitorAddToHistoryResult
		err = workflow.ExecuteActivity(ctx, w.activities.MonitorAddToHistory, historyParam).Get(ctx, &historyResult)
		if err != nil {
			return models.MonitorUnknown, err
		}
	}

	return models.MonitorSuccess, nil
}
