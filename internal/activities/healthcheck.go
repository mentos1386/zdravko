package activities

import (
	"context"
	"log"
	"net/http"
)

type HealtcheckHttpParam struct {
	Url    string
	Method string
}

type HealthcheckHttpResult struct {
	StatusCode int
}

func HealthcheckHttp(ctx context.Context, param HealtcheckHttpParam) (*HealthcheckHttpResult, error) {
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

	return &HealthcheckHttpResult{StatusCode: response.StatusCode}, nil
}

type HealtcheckHttpAddToHistoryParam struct {
	Id         string
	Success    bool
	StatusCode int
}

type HealthcheckHttpAddToHistoryResult struct {
}

func HealthcheckHttpWriteResult(ctx context.Context, param HealtcheckHttpParam) (*HealthcheckHttpResult, error) {
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

	return &HealthcheckHttpResult{StatusCode: response.StatusCode}, nil
}
