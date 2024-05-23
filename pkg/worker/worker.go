package worker

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mentos1386/zdravko/internal/config"
	"github.com/mentos1386/zdravko/internal/temporal"
	"github.com/mentos1386/zdravko/internal/worker/activities"
	"github.com/mentos1386/zdravko/pkg/api"
	"github.com/mentos1386/zdravko/pkg/retry"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

type ConnectionConfig struct {
	Endpoint string `json:"endpoint"`
	Group    string `json:"group"`
}

func getConnectionConfig(token string, apiUrl string) (*ConnectionConfig, error) {
	req, err := api.NewRequest(http.MethodGet, apiUrl+"/api/v1/workers/connect", token, nil)
	if err != nil {
		return nil, err
	}

	return retry.Retry(10, 3*time.Second, func() (*ConnectionConfig, error) {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to API")
		}

		if res.StatusCode == http.StatusUnauthorized {
			panic("WORKER_GROUP_TOKEN is invalid. Either it expired or the worker was removed!")
		}

		if res.StatusCode != http.StatusOK {
			return nil, errors.Errorf("unexpected status code: %d", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}

		config := ConnectionConfig{}
		err = json.Unmarshal(body, &config)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal connection config")
		}

		return &config, nil
	})
}

type Worker struct {
	worker worker.Worker
	cfg    *config.WorkerConfig
	logger *slog.Logger
}

func NewWorker(cfg *config.WorkerConfig) (*Worker, error) {
	return &Worker{
		cfg:    cfg,
		logger: slog.Default().WithGroup("worker"),
	}, nil
}

func (w *Worker) Name() string {
	return "Temporal Worker"
}

func (w *Worker) Start() error {
	config, err := getConnectionConfig(w.cfg.Token, w.cfg.ApiUrl)
	if err != nil {
		return err
	}

	w.logger.Info("Worker Starting", "group", config.Group)

	temporalClient, err := temporal.ConnectWorkerToTemporal(w.cfg.Token, config.Endpoint)
	if err != nil {
		return err
	}

	// Create a new Worker
	w.worker = worker.New(temporalClient, config.Group, worker.Options{})

	workerActivities := activities.NewActivities(w.cfg, w.logger)

	// Register Activities
	w.worker.RegisterActivityWithOptions(workerActivities.Check, activity.RegisterOptions{Name: temporal.ActivityCheckName})

	return w.worker.Run(worker.InterruptCh())
}

func (w *Worker) Stop() error {
	w.worker.Stop()
	return nil
}
