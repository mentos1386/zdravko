package worker

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/temporal"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"code.tjo.space/mentos1386/zdravko/pkg/api"
	"code.tjo.space/mentos1386/zdravko/pkg/retry"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/worker"
)

type ConnectionConfig struct {
	Endpoint string `json:"endpoint"`
	Slug     string `json:"slug"`
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
			panic("WORKER_TOKEN is invalid. Either it expired or the worker was removed!")
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
}

func NewWorker(cfg *config.WorkerConfig) (*Worker, error) {
	return &Worker{
		cfg: cfg,
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

	log.Println("Worker slug:", config.Slug)
	log.Println("Worker group:", config.Group)

	temporalClient, err := temporal.ConnectWorkerToTemporal(w.cfg.Token, config.Endpoint, config.Slug)
	if err != nil {
		return err
	}

	// Create a new Worker
	w.worker = worker.New(temporalClient, config.Group, worker.Options{})

	workerActivities := activities.NewActivities(w.cfg)
	workerWorkflows := workflows.NewWorkflows(workerActivities)

	// Register Workflows
	w.worker.RegisterWorkflow(workerWorkflows.HealthcheckWorkflowDefinition)

	// Register Activities
	w.worker.RegisterActivity(workerActivities.Healthcheck)
	w.worker.RegisterActivity(workerActivities.HealthcheckAddToHistory)

	return w.worker.Run(worker.InterruptCh())
}

func (w *Worker) Stop() error {
	w.worker.Stop()
	return nil
}
