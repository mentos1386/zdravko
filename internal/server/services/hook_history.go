package services

import (
	"context"
	"fmt"
	"time"

	"github.com/mentos1386/zdravko/internal/temporal"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type HookHistory struct {
	HookId          string
	Status          string
	Duration        time.Duration
	StartTime       time.Time
	EndTime         time.Time
	WorkerGroupName string
	Note            string
}

func GetLastNHookHistory(ctx context.Context, t client.Client, n int32) ([]*HookHistory, error) {
	var hookHistory []*HookHistory

	response, err := t.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: n,
	})
	if err != nil {
		return hookHistory, err
	}

	executions := response.GetExecutions()

	for _, execution := range executions {
		scheduleId := string(execution.GetSearchAttributes().GetIndexedFields()["TemporalScheduledById"].Data)

		// Remove the quotes around the hookId and the prefix.
		hookId := scheduleId[len("\"hook-") : len(scheduleId)-1]

		var result temporal.WorkflowHookResult
		if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			workflowRun := t.GetWorkflow(ctx, execution.GetExecution().GetWorkflowId(), execution.GetExecution().GetRunId())
			err := workflowRun.Get(ctx, &result)
			if err != nil {
				return nil, err
			}
		}

		hookHistory = append(hookHistory, &HookHistory{
			HookId:          hookId,
			Duration:        execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			StartTime:       execution.StartTime.AsTime(),
			EndTime:         execution.CloseTime.AsTime(),
			Status:          execution.Status.String(),
			WorkerGroupName: execution.GetTaskQueue(),
			Note:            result.Note,
		})
	}

	return hookHistory, nil
}

func GetHookHistoryForHook(ctx context.Context, t client.Client, hookId string) ([]*HookHistory, error) {
	var hookHistory []*HookHistory

	response, err := t.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: 10,
		Query:    fmt.Sprintf(`TemporalScheduledById = "%s"`, getScheduleId(hookId)),
	})
	if err != nil {
		return hookHistory, err
	}

	executions := response.GetExecutions()

	for _, execution := range executions {
		scheduleId := string(execution.GetSearchAttributes().GetIndexedFields()["TemporalScheduledById"].Data)

		// Remove the quotes around the hookId and the prefix.
		hookId := scheduleId[len("\"hook-") : len(scheduleId)-1]

		var result temporal.WorkflowHookResult
		if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			workflowRun := t.GetWorkflow(ctx, execution.GetExecution().GetWorkflowId(), execution.GetExecution().GetRunId())
			err := workflowRun.Get(ctx, &result)
			if err != nil {
				return nil, err
			}
		}

		hookHistory = append(hookHistory, &HookHistory{
			HookId:          hookId,
			Duration:        execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			StartTime:       execution.StartTime.AsTime(),
			EndTime:         execution.CloseTime.AsTime(),
			Status:          execution.Status.String(),
			WorkerGroupName: execution.GetTaskQueue(),
			Note:            result.Note,
		})
	}

	return hookHistory, nil
}
