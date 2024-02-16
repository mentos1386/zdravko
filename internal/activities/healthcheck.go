package activities

import (
	"context"
	"log"
	"net/http"
)

type HealtcheckHttpActivityParam struct {
	Url    string
	Method string
}

type HealthcheckHttpActivityResult struct {
	StatusCode int
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

	log.Printf("HealthcheckHttpActivityDefinition produced statuscode %d for url %s", response.StatusCode, param.Url)

	return &HealthcheckHttpActivityResult{StatusCode: response.StatusCode}, nil
}
