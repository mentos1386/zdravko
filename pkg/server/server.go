package server

import (
	"context"
	"log"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/handlers"
	"code.tjo.space/mentos1386/zdravko/internal/temporal"
	"code.tjo.space/mentos1386/zdravko/web/static"
	"code.tjo.space/mentos1386/zdravko/web/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	server *http.Server
	cfg    *config.ServerConfig
}

func NewServer(cfg *config.ServerConfig) (*Server, error) {
	return &Server{
		cfg: cfg,
	}, nil
}

func (s *Server) Name() string {
	return "HTTP WEB and API Server"
}

func (s *Server) Start() error {
	e := echo.New()
	e.Renderer = templates.NewTemplates()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, query, err := internal.ConnectToDatabase(s.cfg.DatabasePath)
	if err != nil {
		return err
	}
	log.Println("Connected to database")

	temporalClient, err := temporal.ConnectServerToTemporal(s.cfg)
	if err != nil {
		return err
	}
	log.Println("Connected to Temporal")

	h := handlers.NewBaseHandler(db, query, temporalClient, s.cfg)

	// Health
	e.GET("/health", func(c echo.Context) error {
		d, err := db.DB()
		if err != nil {
			return err
		}
		err = d.Ping()
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Server static files
	stat := e.Group("/static")
	stat.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(static.Static),
	}))

	// Public
	e.GET("", h.Index)

	// Settings
	settings := e.Group("/settings")
	settings.Use(h.Authenticated)
	settings.GET("", h.SettingsOverviewGET)
	settings.GET("/healthchecks", h.SettingsHealthchecksGET)
	settings.GET("/healthchecks/create", h.SettingsHealthchecksCreateGET)
	settings.POST("/healthchecks/create", h.SettingsHealthchecksCreatePOST)
	settings.GET("/healthchecks/:slug", h.SettingsHealthchecksDescribeGET)
	settings.GET("/workers", h.SettingsWorkersGET)
	settings.GET("/workers/create", h.SettingsWorkersCreateGET)
	settings.POST("/workers/create", h.SettingsWorkersCreatePOST)
	settings.GET("/workers/:slug", h.SettingsWorkersDescribeGET)
	settings.GET("/workers/:slug/token", h.SettingsWorkersTokenGET)
	settings.Match([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}, "/temporal*", h.Temporal)

	// OAuth2
	oauth2 := e.Group("/oauth2")
	oauth2.GET("/login", h.OAuth2LoginGET)
	oauth2.GET("/callback", h.OAuth2CallbackGET)
	oauth2.GET("/logout", h.OAuth2LogoutGET, h.Authenticated)

	// API
	apiv1 := e.Group("/api/v1")
	apiv1.Use(h.Authenticated)
	apiv1.GET("/workers/connect", h.ApiV1WorkersConnectGET)
	apiv1.POST("/healthcheck/:slug/history", h.ApiV1HealthchecksHistoryPOST)

	// Error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
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

	return e.Start(":" + s.cfg.Port)
}

func (s *Server) Stop() error {
	ctx := context.Background()
	return s.server.Shutdown(ctx)
}
