package workflows

import (
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"go.temporal.io/sdk/workflow"
)

type HealthcheckHttpWorkflowParam struct {
	Id uint
}

func HealthcheckHttpWorkflowDefinition(ctx workflow.Context, param HealthcheckHttpWorkflowParam) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	activityParam := activities.HealtcheckHttpActivityParam{
		Url:    "https://google.com",
		Method: "GET",
	}

	var result *activities.HealthcheckHttpActivityResult
	err := workflow.ExecuteActivity(ctx, activities.HealthcheckHttpActivityDefinition, activityParam).Get(ctx, &result)

	return err
}
