package worker

import (
	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/temporal"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	worker worker.Worker
	cfg    *config.Config
}

func NewWorker(cfg *config.Config) (*Worker, error) {
	return &Worker{
		cfg: cfg,
	}, nil
}

func (w *Worker) Name() string {
	return "Temporal Worker"
}

func (w *Worker) Start() error {
	temporalClient, err := temporal.ConnectWorkerToTemporal(w.cfg)
	if err != nil {
		return err
	}

	// Create a new Worker
	// TODO: Maybe identify by region or something?
	w.worker = worker.New(temporalClient, "test", worker.Options{})

	// Register Workflows
	w.worker.RegisterWorkflow(workflows.HealthcheckHttpWorkflowDefinition)

	// Register Activities
	w.worker.RegisterActivity(activities.HealthcheckHttp)

	return w.worker.Run(worker.InterruptCh())
}

func (w *Worker) Stop() error {
	w.worker.Stop()
	return nil
}
