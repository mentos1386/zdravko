package activities

import (
	"context"
	"net/http"
)

type HealtcheckHttpActivityParam struct {
	Url    string
	Method string
}

type HealthcheckHttpActivityResult struct {
	Success bool
}

func HealthcheckHttpActivityDefinition(ctx context.Context, param HealtcheckHttpActivityParam) (*HealthcheckHttpActivityResult, error) {
	if param.Method == "" {
		param.Method = "GET"
	}

	var (
		response *http.Response
		err      error
	)

	switch param.Method {
	case "GET":
		response, err = http.Get(param.Url)
	case "POST":
		response, err = http.Post(param.Url, "application/json", nil)
	}

	if err != nil {
		return nil, err
	}

	return &HealthcheckHttpActivityResult{Success: response.StatusCode == 200}, nil
}
