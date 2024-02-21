package k6

import (
	"context"
	"testing"

	"go.k6.io/k6/cmd/state"
)

func TestK6(t *testing.T) {
	ctx := context.Background()

	state := state.NewGlobalState(ctx)

	script := `
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function () {
  http.get('http://test.k6.io');
  sleep(1);
}
`

	execution := NewExecution(state, script)

	err := execution.Start(ctx)
	if err != nil {
		t.Errorf("Error starting execution: %v", err)
	}
}
