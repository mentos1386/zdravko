package workflows

import (
	"fmt"
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"go.temporal.io/sdk/workflow"
)

type HealthcheckHttpWorkflowParam struct {
	Url    string
	Method string
}

func HealthcheckHttpWorkflowDefinition(ctx workflow.Context, param HealthcheckHttpWorkflowParam) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	activityParam := activities.HealtcheckHttpParam{
		Url:    param.Url,
		Method: param.Method,
	}

	var result *activities.HealthcheckHttpResult
	err := workflow.ExecuteActivity(ctx, activities.HealthcheckHttp, activityParam).Get(ctx, &result)
	if err != nil {
		return err
	}

	if result.StatusCode != 200 {
		return fmt.Errorf("HealthcheckHttpActivityDefinition produced statuscode %d for url %s", result.StatusCode, param.Url)
	}

	return nil
}
