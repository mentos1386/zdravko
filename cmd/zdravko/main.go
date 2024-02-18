package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/pkg/server"
	"code.tjo.space/mentos1386/zdravko/pkg/temporal"
	"code.tjo.space/mentos1386/zdravko/pkg/worker"
)

type StartableAndStoppable interface {
	Name() string
	Start() error
	Stop() error
}

func main() {
	var startServer bool
	var startWorker bool
	var startTemporal bool

	flag.BoolVar(&startServer, "server", false, "Start the server")
	flag.BoolVar(&startWorker, "worker", false, "Start the worker")
	flag.BoolVar(&startTemporal, "temporal", false, "Start the temporal")

	flag.Parse()

	println("Starting zdravko...")
	println("Server:   ", startServer)
	println("Worker:   ", startWorker)
	println("Temporal: ", startTemporal)

	if !startServer && !startWorker && !startTemporal {
		log.Fatal("At least one of the following must be set: --server, --worker, --temporal")
	}

	cfg := config.NewConfig()

	var servers [3]StartableAndStoppable
	var wg sync.WaitGroup

	if startTemporal {
		log.Println("Setting up Temporal")
		temporal, err := temporal.NewTemporal(cfg)
		if err != nil {
			log.Fatalf("Unable to create temporal: %v", err)
		}
		servers[0] = temporal
	}

	if startServer {
		log.Println("Setting up Server")
		server, err := server.NewServer(cfg)
		if err != nil {
			log.Fatalf("Unable to create server: %v", err)
		}
		servers[1] = server
	}

	if startWorker {
		log.Println("Setting up Worker")
		worker, err := worker.NewWorker(cfg)
		if err != nil {
			log.Fatalf("Unable to create worker: %v", err)
		}
		servers[2] = worker
	}

	for _, s := range servers {
		srv := s
		if srv == nil {
			continue
		}

		println("Starting", srv.Name())
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := srv.Start()
			if err != nil {
				log.Fatalf("Unable to start server %s: %v", srv.Name(), err)
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("Received signal: %v", sig)
			for _, srv := range servers {
				if srv == nil {
					continue
				}

				println("Stopping", srv.Name())
				err := srv.Stop()
				if err != nil {
					log.Printf("Unable to stop server %s: %v", srv.Name(), err)
				}
			}
		}
	}()

	wg.Wait()
}
