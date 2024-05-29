package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mentos1386/zdravko/internal/config"
	"github.com/mentos1386/zdravko/pkg/server"
	"github.com/mentos1386/zdravko/pkg/temporal"
	"github.com/mentos1386/zdravko/pkg/worker"
)

type StartableAndStoppable interface {
	Name() string
	Start() error
	Stop() error
}

func setupLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)
}

func main() {
	setupLogger()

	var startServer bool
	var startWorker bool
	var startTemporal bool

	flag.BoolVar(&startServer, "server", false, "Start the server")
	flag.BoolVar(&startWorker, "worker", false, "Start the worker")
	flag.BoolVar(&startTemporal, "temporal", false, "Start the temporal")
	flag.Parse()

	slog.Info("Starting zdravko...", "server", startServer, "worker", startWorker, "temporal", startTemporal)

	if !startServer && !startWorker && !startTemporal {
		slog.Error("At least one of the following must be set: --server, --worker, --temporal")
		os.Exit(1)
	}

	var servers [3]StartableAndStoppable
	var wg sync.WaitGroup

	if startTemporal {
		slog.Info("Setting up Temporal")
		cfg := config.NewTemporalConfig()
		temporal, err := temporal.NewTemporal(cfg)
		if err != nil {
			slog.Error("Unable to create temporal", "error", err)
			os.Exit(1)
		}
		servers[0] = temporal
	}

	if startServer {
		slog.Info("Setting up Server")
		cfg := config.NewServerConfig()
		server, err := server.NewServer(cfg)
		if err != nil {
			slog.Error("Unable to create server", "error", err)
			os.Exit(1)
		}
		servers[1] = server
	}

	if startWorker {
		slog.Info("Setting up Worker")
		cfg := config.NewWorkerConfig()
		worker, err := worker.NewWorker(cfg)
		if err != nil {
			slog.Error("Unable to create worker", "error", err)
			os.Exit(1)
		}
		servers[2] = worker
	}

	for _, s := range servers {
		srv := s
		if srv == nil {
			continue
		}

		slog.Info("Starting", "name", srv.Name())
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := srv.Start()
			if err != nil {
				slog.Error("Unable to start", "name", srv.Name(), "error", err)
				os.Exit(1)
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			slog.Info("Received signal", "signal", sig)
			for _, srv := range servers {
				if srv == nil {
					continue
				}

				slog.Info("Stopping", "name", srv.Name())
				err := srv.Stop()
				if err != nil {
					slog.Error("Unable to stop server", "name", srv.Name(), "error", err)
				}
			}
		}
	}()

	wg.Wait()
}
