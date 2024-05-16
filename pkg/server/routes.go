package server

import (
	"log/slog"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/handlers"
	"code.tjo.space/mentos1386/zdravko/internal/kv"
	"code.tjo.space/mentos1386/zdravko/web/static"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.temporal.io/sdk/client"
)

func Routes(
	e *echo.Echo,
	sqlDb *sqlx.DB,
	kvStore kv.KeyValueStore,
	temporalClient client.Client,
	cfg *config.ServerConfig,
	logger *slog.Logger,
) {
	h := handlers.NewBaseHandler(sqlDb, kvStore, temporalClient, cfg, logger)

	// Health
	e.GET("/health", func(c echo.Context) error {
		err := sqlDb.Ping()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Server static files
	stat := e.Group("/static")
	stat.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "public, max-age=60")
			return next(c)
		}
	})
	stat.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(static.Static),
	}))

	// Public
	e.GET("", h.Index)
	e.GET("/incidents", h.Incidents)

	// Settings
	settings := e.Group("/settings")
	settings.Use(h.Authenticated)
	settings.GET("", h.SettingsOverviewGET)

	settings.GET("/triggers", h.SettingsTriggersGET)
	settings.GET("/triggers/:id", h.SettingsTriggersDescribeGET)
	settings.POST("/triggers/:id", h.SettingsTriggersDescribePOST)
	settings.GET("/triggers/create", h.SettingsTriggersCreateGET)
	settings.POST("/triggers/create", h.SettingsTriggersCreatePOST)
	settings.GET("/triggers/:id/delete", h.SettingsTriggersDescribeDELETE)
	settings.GET("/triggers/:id/disable", h.SettingsTriggersDisableGET)
	settings.GET("/triggers/:id/enable", h.SettingsTriggersEnableGET)

	settings.GET("/targets", h.SettingsTargetsGET)
	//settings.GET("/targets/create", h.SettingsTargetsCreateGET)
	//settings.POST("/targets/create", h.SettingsTargetsCreatePOST)
	//settings.GET("/targets/:id", h.SettingsTargetsDescribeGET)
	//settings.POST("/targets/:id", h.SettingsTargetsDescribePOST)
	//settings.GET("/targets/:id/delete", h.SettingsTargetsDescribeDELETE)
	//settings.GET("/targets/:id/disable", h.SettingsTargetsDisableGET)
	//settings.GET("/targets/:id/enable", h.SettingsTargetsEnableGET)

	settings.GET("/incidents", h.SettingsIncidentsGET)

	settings.GET("/checks", h.SettingsChecksGET)
	settings.GET("/checks/create", h.SettingsChecksCreateGET)
	settings.POST("/checks/create", h.SettingsChecksCreatePOST)
	settings.GET("/checks/:id", h.SettingsChecksDescribeGET)
	settings.POST("/checks/:id", h.SettingsChecksDescribePOST)
	settings.GET("/checks/:id/delete", h.SettingsChecksDescribeDELETE)
	settings.GET("/checks/:id/disable", h.SettingsChecksDisableGET)
	settings.GET("/checks/:id/enable", h.SettingsChecksEnableGET)

	settings.GET("/notifications", h.SettingsNotificationsGET)

	settings.GET("/worker-groups", h.SettingsWorkerGroupsGET)
	settings.GET("/worker-groups/create", h.SettingsWorkerGroupsCreateGET)
	settings.POST("/worker-groups/create", h.SettingsWorkerGroupsCreatePOST)
	settings.GET("/worker-groups/:id", h.SettingsWorkerGroupsDescribeGET)
	settings.GET("/worker-groups/:id/delete", h.SettingsWorkerGroupsDescribeDELETE)

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
	apiv1.POST("/checks/:id/history", h.ApiV1ChecksHistoryPOST)

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
}
