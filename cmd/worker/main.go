package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Initialize a Temporal Client
	// Specify the Namespace in the Client options
	clientOptions := client.Options{
		Namespace: "default",
	}
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create a Temporal Client", err)
	}
	defer temporalClient.Close()
	// Create a new Worker
	yourWorker := worker.New(temporalClient, "default-boilerplate-task-queue-local", worker.Options{})
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
