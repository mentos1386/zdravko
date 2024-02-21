package activities

import (
	"context"
	"log/slog"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/pkg/k6"
)

type HealtcheckParam struct {
	Script string
}

type HealthcheckResult struct {
	StatusCode int
}

func Healthcheck(ctx context.Context, param HealtcheckParam) (*HealthcheckResult, error) {

	statusCode := http.StatusOK // FIXME

	execution := k6.NewExecution(slog.Default(), param.Script)

	err := execution.Run(ctx)
	if err != nil {
		return nil, err
	}

	return &HealthcheckResult{StatusCode: statusCode}, nil
}

type HealtcheckAddToHistoryParam struct {
	Id         string
	Success    bool
	StatusCode int
}

type HealthcheckAddToHistoryResult struct {
}

func HealthcheckAddToHistory(ctx context.Context, param HealtcheckAddToHistoryParam) (*HealthcheckAddToHistoryResult, error) {

	return &HealthcheckAddToHistoryResult{}, nil
}
