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

type CheckHistory struct {
	CheckId         string
	Status          string
	Duration        time.Duration
	StartTime       time.Time
	EndTime         time.Time
	WorkerGroupName string
	Note            string
}

func GetLastNCheckHistory(ctx context.Context, t client.Client, n int32) ([]*CheckHistory, error) {
	var checkHistory []*CheckHistory

	response, err := t.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: n,
	})
	if err != nil {
		return checkHistory, err
	}

	executions := response.GetExecutions()

	for _, execution := range executions {
		scheduleId := string(execution.GetSearchAttributes().GetIndexedFields()["TemporalScheduledById"].Data)

		// Remove the quotes around the checkId and the prefix.
		checkId := scheduleId[len("\"check-") : len(scheduleId)-1]

		var result temporal.WorkflowCheckResult
		if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			workflowRun := t.GetWorkflow(ctx, execution.GetExecution().GetWorkflowId(), execution.GetExecution().GetRunId())
			err := workflowRun.Get(ctx, &result)
			if err != nil {
				return nil, err
			}
		}

		checkHistory = append(checkHistory, &CheckHistory{
			CheckId:         checkId,
			Duration:        execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			StartTime:       execution.StartTime.AsTime(),
			EndTime:         execution.CloseTime.AsTime(),
			Status:          execution.Status.String(),
			WorkerGroupName: execution.GetTaskQueue(),
			Note:            result.Note,
		})
	}

	return checkHistory, nil
}

func GetCheckHistoryForCheck(ctx context.Context, t client.Client, checkId string) ([]*CheckHistory, error) {
	var checkHistory []*CheckHistory

	response, err := t.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: 10,
		Query:    fmt.Sprintf(`TemporalScheduledById = "%s"`, getScheduleId(checkId)),
	})
	if err != nil {
		return checkHistory, err
	}

	executions := response.GetExecutions()

	for _, execution := range executions {
		scheduleId := string(execution.GetSearchAttributes().GetIndexedFields()["TemporalScheduledById"].Data)

		// Remove the quotes around the checkId and the prefix.
		checkId := scheduleId[len("\"check-") : len(scheduleId)-1]

		var result temporal.WorkflowCheckResult
		if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			workflowRun := t.GetWorkflow(ctx, execution.GetExecution().GetWorkflowId(), execution.GetExecution().GetRunId())
			err := workflowRun.Get(ctx, &result)
			if err != nil {
				return nil, err
			}
		}

		checkHistory = append(checkHistory, &CheckHistory{
			CheckId:         checkId,
			Duration:        execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			StartTime:       execution.StartTime.AsTime(),
			EndTime:         execution.CloseTime.AsTime(),
			Status:          execution.Status.String(),
			WorkerGroupName: execution.GetTaskQueue(),
			Note:            result.Note,
		})
	}

	return checkHistory, nil
}
