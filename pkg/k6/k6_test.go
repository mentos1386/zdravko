package k6

import (
	"context"
	"log/slog"
	"testing"
)

func TestK6Success(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	script := `
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '5s',
};

export default function () {
  http.get('https://test.k6.io');
  sleep(1);
}
`

	execution := NewExecution(logger, script)

	err := execution.Run(ctx)
	if err != nil {
		t.Errorf("Error starting execution: %v", err)
	}
}

func TestK6Fail(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	script := `
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  thresholds: {
    http_req_failed: ['rate<0.01'], // http errors should be less than 1%
  },
};

export default function () {
  http.get('https://fail.broken.example');
  sleep(1);
}
`

	execution := NewExecution(logger, script)

	err := execution.Run(ctx)
	if err != nil {
		t.Errorf("Error starting execution: %v", err)
	}
}
