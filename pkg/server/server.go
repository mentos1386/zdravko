package server

import (
	"context"
	"log/slog"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/database"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/handlers"
	"code.tjo.space/mentos1386/zdravko/internal/temporal"
	"code.tjo.space/mentos1386/zdravko/web/static"
	"code.tjo.space/mentos1386/zdravko/web/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.ServerConfig
	logger *slog.Logger
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
	s.echo.Renderer = templates.NewTemplates()
	//s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.Recover())

	db, err := database.ConnectToDatabase(s.logger, s.cfg.DatabasePath)
	if err != nil {
		return err
	}

	temporalClient, err := temporal.ConnectServerToTemporal(s.logger, s.cfg)
	if err != nil {
		return err
	}

	h := handlers.NewBaseHandler(db, temporalClient, s.cfg, s.logger)

	// Health
	s.echo.GET("/health", func(c echo.Context) error {
		err = db.Ping()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Server static files
	stat := s.echo.Group("/static")
	stat.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(static.Static),
	}))

	// Public
	s.echo.GET("", h.Index)
	s.echo.GET("/incidents", h.Incidents)

	// Settings
	settings := s.echo.Group("/settings")
	settings.Use(h.Authenticated)
	settings.GET("", h.SettingsOverviewGET)
	settings.GET("/monitors", h.SettingsMonitorsGET)
	settings.GET("/monitors/create", h.SettingsMonitorsCreateGET)
	settings.POST("/monitors/create", h.SettingsMonitorsCreatePOST)
	settings.GET("/monitors/:slug", h.SettingsMonitorsDescribeGET)
	settings.POST("/monitors/:slug", h.SettingsMonitorsDescribePOST)
	settings.GET("/worker-groups", h.SettingsWorkerGroupsGET)
	settings.GET("/worker-groups/create", h.SettingsWorkerGroupsCreateGET)
	settings.POST("/worker-groups/create", h.SettingsWorkerGroupsCreatePOST)
	settings.GET("/worker-groups/:slug", h.SettingsWorkerGroupsDescribeGET)
	settings.Match([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}, "/temporal*", h.Temporal)

	// OAuth2
	oauth2 := s.echo.Group("/oauth2")
	oauth2.GET("/login", h.OAuth2LoginGET)
	oauth2.GET("/callback", h.OAuth2CallbackGET)
	oauth2.GET("/logout", h.OAuth2LogoutGET, h.Authenticated)

	// API
	apiv1 := s.echo.Group("/api/v1")
	apiv1.Use(h.Authenticated)
	apiv1.GET("/workers/connect", h.ApiV1WorkersConnectGET)
	apiv1.POST("/monitors/:slug/history", h.ApiV1MonitorsHistoryPOST)

	// Error handler
	s.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		if code == http.StatusNotFound {
			_ = h.Error404(c)
			return
		}
		_ = c.String(code, err.Error())
	}

	return s.echo.Start(":" + s.cfg.Port)
}

func (s *Server) Stop() error {
	ctx := context.Background()
	return s.echo.Shutdown(ctx)
}
