package services

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type CheckHistory struct {
	CheckId  string
	Status   string
	Duration time.Duration
}

func GetLastNCheckHistory(ctx context.Context, temporal client.Client, n int32) ([]*CheckHistory, error) {
	var checkHistory []*CheckHistory

	response, err := temporal.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: n,
	})
	if err != nil {
		return checkHistory, err
	}

	executions := response.GetExecutions()

	for _, execution := range executions {
		scheduleId := string(execution.GetSearchAttributes().GetIndexedFields()["TemporalScheduledById"].Data)
		checkId := scheduleId[len("check-"):]
		checkHistory = append(checkHistory, &CheckHistory{
			CheckId:  checkId,
			Duration: execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			Status:   execution.Status.String(),
		})
	}

	return checkHistory, nil
}

func GetCheckHistoryForCheck(ctx context.Context, temporal client.Client, checkId string) ([]*CheckHistory, error) {
	var checkHistory []*CheckHistory

	response, err := temporal.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: 10,
		Query:    fmt.Sprintf(`TemporalScheduledById = "%s"`, getScheduleId(checkId)),
	})
	if err != nil {
		return checkHistory, err
	}

	executions := response.GetExecutions()

	for _, execution := range executions {
		scheduleId := string(execution.GetSearchAttributes().GetIndexedFields()["TemporalScheduledById"].Data)
		checkId := scheduleId[len("check-"):]
		checkHistory = append(checkHistory, &CheckHistory{
			CheckId:  checkId,
			Duration: execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			Status:   execution.Status.String(),
		})
	}

	return checkHistory, nil
}
