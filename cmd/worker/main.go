package main

import (
	"log"

	"code.tjo.space/mentos1386/zdravko/internal/activities"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	config := config.NewConfig()

	// Initialize a Temporal Client
	// Specify the Namespace in the Client options
	clientOptions := client.Options{
		HostPort:  config.Temporal.ServerHost,
		Namespace: "default",
	}
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create a Temporal Client", err)
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
