package server

import (
	"context"
	"log/slog"

	"code.tjo.space/mentos1386/zdravko/database"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/temporal"
	"code.tjo.space/mentos1386/zdravko/web/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.ServerConfig
	logger *slog.Logger

	worker *Worker
}

func NewServer(cfg *config.ServerConfig) (*Server, error) {
	return &Server{
		cfg:    cfg,
		echo:   echo.New(),
		logger: slog.Default().WithGroup("server"),
	}, nil
}

func (s *Server) Name() string {
	return "HTTP WEB and API Server"
}

func (s *Server) Start() error {
	db, err := database.ConnectToDatabase(s.logger, s.cfg.DatabasePath)
	if err != nil {
		return err
	}

	temporalClient, err := temporal.ConnectServerToTemporal(s.logger, s.cfg)
	if err != nil {
		return err
	}

	s.worker = NewWorker(temporalClient, s.cfg)

	s.echo.Renderer = templates.NewTemplates()
	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Secure())
	Routes(s.echo, db, temporalClient, s.cfg, s.logger)

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
