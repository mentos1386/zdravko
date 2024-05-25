package services

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type TriggerHistory struct {
	TriggerId string
	Status    string
	Duration  time.Duration
}

func GetLastNTriggerHistory(ctx context.Context, temporal client.Client, n int32) ([]*TriggerHistory, error) {
	var checkHistory []*TriggerHistory

	response, err := temporal.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: n,
	})
	if err != nil {
		return checkHistory, err
	}

	executions := response.GetExecutions()

	for _, execution := range executions {
		scheduleId := string(execution.GetSearchAttributes().GetIndexedFields()["TemporalScheduledById"].Data)
		checkId := scheduleId[len("trigger-"):]
		checkHistory = append(checkHistory, &TriggerHistory{
			TriggerId: checkId,
			Duration:  execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			Status:    execution.Status.String(),
		})
	}

	return checkHistory, nil
}

func GetTriggerHistoryForTrigger(ctx context.Context, temporal client.Client, checkId string) ([]*TriggerHistory, error) {
	var checkHistory []*TriggerHistory

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
		checkHistory = append(checkHistory, &TriggerHistory{
			TriggerId: checkId,
			Duration:  execution.CloseTime.AsTime().Sub(execution.StartTime.AsTime()),
			Status:    execution.Status.String(),
		})
	}

	return checkHistory, nil
}
