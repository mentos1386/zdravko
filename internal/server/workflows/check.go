package workflows

import (
	"log/slog"
	"sort"
	"time"

	"github.com/mentos1386/zdravko/internal/temporal"
	"go.temporal.io/sdk/workflow"
)

func (w *Workflows) CheckWorkflowDefinition(ctx workflow.Context, param temporal.WorkflowCheckParam) error {
	workerGroupIds := param.WorkerGroupIds
	sort.Strings(workerGroupIds)

	targetsFilterResult := temporal.ActivityTargetsFilterResult{}
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 60 * time.Second,
			TaskQueue:           temporal.TEMPORAL_SERVER_QUEUE,
		}),
		temporal.ActivityTargetsFilterName,
		temporal.ActivityTargetsFilterParam{
			Filter: param.Filter,
		},
	).Get(ctx, &targetsFilterResult)
	if err != nil {
		return err
	}

	for _, target := range targetsFilterResult.Targets {
		for _, workerGroupId := range workerGroupIds {
			var checkResult *temporal.ActivityCheckResult
			err := workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
					StartToCloseTimeout: 60 * time.Second,
					TaskQueue:           workerGroupId,
				}),
				temporal.ActivityCheckName,
				temporal.ActivityCheckParam{
					Script: param.Script,
					Target: target,
				},
			).Get(ctx, &checkResult)
			if err != nil {
				return err
			}

			status := temporal.AddTargetHistoryStatusFailure
			if checkResult.Success {
				status = temporal.AddTargetHistoryStatusSuccess
			}

			var addTargetHistoryResult *temporal.ActivityAddTargetHistoryResult
			err = workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
					StartToCloseTimeout: 60 * time.Second,
					TaskQueue:           temporal.TEMPORAL_SERVER_QUEUE,
				}),
				temporal.ActivityAddTargetHistoryName,
				&temporal.ActivityAddTargetHistoryParam{
					Target: target,
					Status: status,
					Note:   checkResult.Note,
				},
			).Get(ctx, &addTargetHistoryResult)
			if err != nil {
				return err
			}

			slog.Info("Check %s status: %s", param.CheckId, status)
		}
	}

	return nil
}
