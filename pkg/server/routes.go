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
	settings.GET("/incidents", h.SettingsIncidentsGET)
	settings.GET("/monitors", h.SettingsMonitorsGET)
	settings.GET("/monitors/create", h.SettingsMonitorsCreateGET)
	settings.POST("/monitors/create", h.SettingsMonitorsCreatePOST)
	settings.GET("/monitors/:id", h.SettingsMonitorsDescribeGET)
	settings.POST("/monitors/:id", h.SettingsMonitorsDescribePOST)
	settings.GET("/monitors/:id/delete", h.SettingsMonitorsDescribeDELETE)
	settings.GET("/monitors/:id/disable", h.SettingsMonitorsDisableGET)
	settings.GET("/monitors/:id/enable", h.SettingsMonitorsEnableGET)
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
	apiv1.POST("/monitors/:id/history", h.ApiV1MonitorsHistoryPOST)

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
