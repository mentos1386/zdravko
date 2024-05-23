package workflows

import (
	"log/slog"
	"sort"
	"time"

	"github.com/mentos1386/zdravko/internal/temporal"
	"github.com/mentos1386/zdravko/pkg/api"
	"go.temporal.io/sdk/workflow"
)

func (w *Workflows) CheckWorkflowDefinition(ctx workflow.Context, param temporal.WorkflowCheckParam) (api.CheckStatus, error) {
	workerGroupIds := param.WorkerGroupIds
	sort.Strings(workerGroupIds)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           temporal.TEMPORAL_SERVER_QUEUE,
	})
	targetsFilterParam := temporal.ActivityTargetsFilterParam{
		Filter: param.Filter,
	}
	targetsFilterResult := temporal.ActivityTargetsFilterResult{}
	err := workflow.ExecuteActivity(ctx, temporal.ActivityTargetsFilterName, targetsFilterParam).Get(ctx, &targetsFilterResult)
	if err != nil {
		return api.CheckStatusUnknown, err
	}

	for _, target := range targetsFilterResult.Targets {
		for _, workerGroupId := range workerGroupIds {
			ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 60 * time.Second,
				TaskQueue:           workerGroupId,
			})

			heatlcheckParam := temporal.ActivityCheckParam{
				Script: param.Script,
				Target: target,
			}

			var checkResult *temporal.ActivityCheckResult
			err := workflow.ExecuteActivity(ctx, temporal.ActivityCheckName, heatlcheckParam).Get(ctx, &checkResult)
			if err != nil {
				return api.CheckStatusUnknown, err
			}

			status := api.CheckStatusFailure
			if checkResult.Success {
				status = api.CheckStatusSuccess
			}

			slog.Info("Check %s status: %s", param.CheckId, status)
		}
	}

	return api.CheckStatusSuccess, nil
}
