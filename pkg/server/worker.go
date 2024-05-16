package server

import (
	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	worker worker.Worker
}

func NewWorker(temporalClient client.Client, cfg *config.ServerConfig) *Worker {
	w := worker.New(temporalClient, "default", worker.Options{})

	workerActivities := activities.NewActivities(&config.WorkerConfig{})

	workerWorkflows := workflows.NewWorkflows(workerActivities)

	// Register Workflows
	w.RegisterWorkflow(workerWorkflows.CheckWorkflowDefinition)

	return &Worker{
		worker: w,
	}
}

func (w *Worker) Start() error {
	return w.worker.Run(worker.InterruptCh())
}

func (w *Worker) Stop() {
	w.worker.Stop()
}
