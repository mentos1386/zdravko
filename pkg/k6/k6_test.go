package k6

import (
	"context"
	"log/slog"
	"testing"
)

func TestK6(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	script := `
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  http.get('https://test.k6.io');
  sleep(1);
}
`

	execution := NewExecution(logger, script)

	err := execution.Start(ctx)
	if err != nil {
		t.Errorf("Error starting execution: %v", err)
	}
}
