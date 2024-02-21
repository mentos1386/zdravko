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
	url := apiUrl + "/api/v1/workers/connect"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	return retry.Retry(10, 3*time.Second, func() (*ConnectionConfig, error) {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to API")
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
	// TODO: Maybe identify by region or something?
	w.worker = worker.New(temporalClient, config.Group, worker.Options{})

	// Register Workflows
	w.worker.RegisterWorkflow(workflows.HealthcheckWorkflowDefinition)

	// Register Activities
	w.worker.RegisterActivity(activities.Healthcheck)

	return w.worker.Run(worker.InterruptCh())
}

func (w *Worker) Stop() error {
	w.worker.Stop()
	return nil
}
