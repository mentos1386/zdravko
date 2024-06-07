package server

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mentos1386/zdravko/database"
	"github.com/mentos1386/zdravko/internal/config"
	"github.com/mentos1386/zdravko/internal/temporal"
	"github.com/mentos1386/zdravko/web/templates"
	"github.com/pkg/errors"
)

type Server struct {
	echo    *echo.Echo
	cfg     *config.ServerConfig
	logger  *slog.Logger
	version string

	worker *Worker
}

func NewServer(version string, cfg *config.ServerConfig) (*Server, error) {
	return &Server{
		cfg:     cfg,
		echo:    echo.New(),
		logger:  slog.Default(),
		version: version,
	}, nil
}

func (s *Server) Name() string {
	return "HTTP WEB and API Server"
}

func (s *Server) Start() error {
	sqliteDb, err := database.ConnectToDatabase(s.logger, s.cfg.SqliteDatabasePath)
	if err != nil {
		return errors.Wrap(err, "failed to connect to sqlite")
	}

	temporalClient, err := temporal.ConnectServerToTemporal(s.logger, s.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to connect to temporal")
	}

	kvStore, err := database.NewBadgerKeyValueStore(s.cfg.KeyValueDatabasePath)
	if err != nil {
		return errors.Wrap(err, "failed to open kv store")
	}

	s.worker = NewWorker(temporalClient, s.cfg, s.logger, sqliteDb, kvStore)

	templates, err := templates.NewTemplates(s.version, s.logger)
	if err != nil {
		return errors.Wrap(err, "failed to create templates")
	}
	s.echo.Renderer = templates

	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Secure())
	Routes(s.echo, sqliteDb, kvStore, temporalClient, s.cfg, s.logger)

	go func() {
		if err := s.worker.Start(); err != nil {
			panic(err)
		}
	}()

	return s.echo.Start(":" + s.cfg.Port)
}

func (s *Server) Stop() error {
	s.worker.Stop()

	ctx := context.Background()
	return s.echo.Shutdown(ctx)
}
