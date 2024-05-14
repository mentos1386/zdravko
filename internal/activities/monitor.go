package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/script"
	"code.tjo.space/mentos1386/zdravko/pkg/api"
	"code.tjo.space/mentos1386/zdravko/pkg/k6"
)

type HealtcheckParam struct {
	Script string
}

type MonitorResult struct {
	Success bool
	Note    string
}

func (a *Activities) Monitor(ctx context.Context, param HealtcheckParam) (*MonitorResult, error) {
	execution := k6.NewExecution(slog.Default(), script.UnescapeString(param.Script))

	result, err := execution.Run(ctx)
	if err != nil {
		return nil, err
	}

	return &MonitorResult{Success: result.Success, Note: result.Note}, nil
}

type HealtcheckAddToHistoryParam struct {
	MonitorId     string
	Status        models.MonitorStatus
	Note          string
	WorkerGroupId string
}

type MonitorAddToHistoryResult struct {
}

func (a *Activities) MonitorAddToHistory(ctx context.Context, param HealtcheckAddToHistoryParam) (*MonitorAddToHistoryResult, error) {
	url := fmt.Sprintf("%s/api/v1/monitors/%s/history", a.config.ApiUrl, param.MonitorId)

	body := api.ApiV1MonitorsHistoryPOSTBody{
		Status:        param.Status,
		Note:          param.Note,
		WorkerGroupId: param.WorkerGroupId,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := api.NewRequest(http.MethodPost, url, a.config.Token, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return &MonitorAddToHistoryResult{}, nil
}
