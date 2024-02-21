package api

import (
	"io"
	"net/http"
)

func NewRequest(method string, urlString string, token string, body io.Reader) (*http.Request, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	if body != nil {
		headers.Add("Content-Type", "application/json")
	}

	req, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	return req, nil
}
