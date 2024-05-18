package workflows

import (
	"sort"
	"time"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"go.temporal.io/sdk/workflow"
)

type CheckWorkflowParam struct {
	Script         string
	CheckId        string
	WorkerGroupIds []string
}

func (w *Workflows) CheckWorkflowDefinition(ctx workflow.Context, param CheckWorkflowParam) (models.CheckStatus, error) {
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

		var checkResult *activities.CheckResult
		err := workflow.ExecuteActivity(ctx, w.activities.Check, heatlcheckParam).Get(ctx, &checkResult)
		if err != nil {
			return models.CheckStatusUnknown, err
		}

		status := models.CheckStatusFailure
		if checkResult.Success {
			status = models.CheckStatusSuccess
		}

		historyParam := activities.HealtcheckAddToHistoryParam{
			CheckId:       param.CheckId,
			Status:        status,
			Note:          checkResult.Note,
			WorkerGroupId: workerGroupId,
		}

		var historyResult *activities.CheckAddToHistoryResult
		err = workflow.ExecuteActivity(ctx, w.activities.CheckAddToHistory, historyParam).Get(ctx, &historyResult)
		if err != nil {
			return models.CheckStatusUnknown, err
		}
	}

	return models.CheckStatusSuccess, nil
}
