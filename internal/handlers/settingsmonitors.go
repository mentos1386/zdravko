package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

type SettingsMonitors struct {
	*Settings
	Monitors       []*models.Monitor
	MonitorsLength int
}

type SettingsMonitor struct {
	*Settings
	Monitor *models.Monitor
}

func (h *BaseHandler) SettingsMonitorsGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	monitors, err := h.query.Monitor.WithContext(context.Background()).Find()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_monitors.tmpl", &SettingsMonitors{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Monitors"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Monitors")},
		),
		Monitors:       monitors,
		MonitorsLength: len(monitors),
	})
}

func (h *BaseHandler) SettingsMonitorsDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("slug")

	monitor, err := services.GetMonitor(context.Background(), h.query, slug)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_monitors_describe.tmpl", &SettingsMonitor{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Monitors"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Monitors"),
				{
					Path:       fmt.Sprintf("/settings/monitors/%s", slug),
					Title:      "Describe",
					Breadcrumb: monitor.Name,
				},
			}),
		Monitor: monitor,
	})
}

func (h *BaseHandler) SettingsMonitorsDescribePOST(c echo.Context) error {
	ctx := context.Background()

	slug := c.Param("slug")

	monitor, err := services.GetMonitor(ctx, h.query, slug)
	if err != nil {
		return err
	}

	update := &models.Monitor{
		Slug:         monitor.Slug,
		Name:         monitor.Name,
		Schedule:     c.FormValue("schedule"),
		WorkerGroups: strings.Split(c.FormValue("workergroups"), " "),
		Script:       c.FormValue("script"),
	}

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(update)
	if err != nil {
		return err
	}

	err = services.UpdateMonitor(
		ctx,
		h.query,
		update,
	)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateMonitorSchedule(ctx, h.temporal, monitor)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/monitors/%s", slug))
}

func (h *BaseHandler) SettingsMonitorsCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_monitors_create.tmpl", NewSettings(
		cc.Principal.User,
		GetPageByTitle(SettingsPages, "Monitors"),
		[]*components.Page{
			GetPageByTitle(SettingsPages, "Monitors"),
			GetPageByTitle(SettingsPages, "Monitors Create"),
		},
	))
}

func (h *BaseHandler) SettingsMonitorsCreatePOST(c echo.Context) error {
	ctx := context.Background()

	monitorHttp := &models.Monitor{
		Name:         c.FormValue("name"),
		Slug:         slug.Make(c.FormValue("name")),
		Schedule:     c.FormValue("schedule"),
		WorkerGroups: strings.Split(c.FormValue("workergroups"), " "),
		Script:       c.FormValue("script"),
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(monitorHttp)
	if err != nil {
		return err
	}

	err = services.CreateMonitor(
		ctx,
		h.query,
		monitorHttp,
	)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateMonitorSchedule(ctx, h.temporal, monitorHttp)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/monitors")
}
