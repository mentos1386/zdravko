package main

import (
	"log"

	"code.tjo.space/mentos1386/zdravko/internal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	config := internal.NewConfig()

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
	yourWorker := worker.New(temporalClient, "default", worker.Options{})

	// Register Workflows
	//yourWorker.RegisterWorkflow(workflows.default)
	// Register Activities
	//yourWorker.RegisterActivity(activities.SSNTraceActivity)
	// Start the the Worker Process
	err = yourWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start the Worker Process", err)
	}
}
