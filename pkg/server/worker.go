package server

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/mentos1386/zdravko/database"
	"github.com/mentos1386/zdravko/internal/config"
	"github.com/mentos1386/zdravko/internal/server/activities"
	"github.com/mentos1386/zdravko/internal/server/workflows"
	"github.com/mentos1386/zdravko/internal/temporal"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	temporalWorker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Worker struct {
	worker temporalWorker.Worker
}

func NewWorker(temporalClient client.Client, cfg *config.ServerConfig, logger *slog.Logger, db *sqlx.DB, kvStore database.KeyValueStore) *Worker {
	worker := temporalWorker.New(temporalClient, temporal.TEMPORAL_SERVER_QUEUE, temporalWorker.Options{})

	a := activities.NewActivities(cfg, logger, db, kvStore)

	w := workflows.NewWorkflows()

	// Register Workflows
	worker.RegisterWorkflowWithOptions(w.CheckWorkflowDefinition, workflow.RegisterOptions{Name: temporal.WorkflowCheckName})

	// Register Activities
	worker.RegisterActivityWithOptions(a.TargetsFilter, activity.RegisterOptions{Name: temporal.ActivityTargetsFilterName})
	worker.RegisterActivityWithOptions(a.ProcessCheckOutcome, activity.RegisterOptions{Name: temporal.ActivityProcessCheckOutcomeName})

	return &Worker{
		worker: worker,
	}
}

func (w *Worker) Start() error {
	return w.worker.Run(temporalWorker.InterruptCh())
}

func (w *Worker) Stop() {
	w.worker.Stop()
}
