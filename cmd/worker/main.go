package main

import (
	"log"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg := config.NewConfig()

	temporalClient, err := internal.ConnectToTemporal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer temporalClient.Close()

	// Create a new Worker
	// TODO: Maybe identify by region or something?
	w := worker.New(temporalClient, "default", worker.Options{})

	// Register Workflows
	w.RegisterWorkflow(workflows.HealthcheckHttpWorkflowDefinition)

	// Register Activities
	w.RegisterActivity(activities.HealthcheckHttpActivityDefinition)

	// Start the the Worker Process
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start the Worker Process", err)
	}
}
