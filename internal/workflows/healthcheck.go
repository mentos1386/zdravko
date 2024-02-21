package workflows

import (
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"go.temporal.io/sdk/workflow"
)

type HealthcheckWorkflowParam struct {
	Script string
}

func HealthcheckWorkflowDefinition(ctx workflow.Context, param HealthcheckWorkflowParam) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	activityParam := activities.HealtcheckParam{
		Script: param.Script,
	}

	var result *activities.HealthcheckResult
	err := workflow.ExecuteActivity(ctx, activities.Healthcheck, activityParam).Get(ctx, &result)
	if err != nil {
		return err
	}

	return nil
}
