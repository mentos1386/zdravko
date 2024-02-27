package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

type CreateMonitor struct {
	Name         string `validate:"required"`
	WorkerGroups string `validate:"required"`
	Schedule     string `validate:"required,cron"`
	Script       string `validate:"required"`
}

type UpdateMonitor struct {
	WorkerGroups string `validate:"required"`
	Schedule     string `validate:"required,cron"`
	Script       string `validate:"required"`
}

type SettingsMonitors struct {
	*Settings
	Monitors       []*models.MonitorWithWorkerGroups
	MonitorsLength int
}

type SettingsMonitor struct {
	*Settings
	Monitor *models.MonitorWithWorkerGroups
	History []*models.MonitorHistory
}

func (h *BaseHandler) SettingsMonitorsGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	monitors, err := services.GetMonitorsWithWorkerGroups(context.Background(), h.db)
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

	monitor, err := services.GetMonitorWithWorkerGroups(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	history, err := services.GetMonitorHistoryForMonitor(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	maxElements := 10
	if len(history) < maxElements {
		maxElements = len(history)
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
		History: history[:maxElements],
	})
}

func (h *BaseHandler) SettingsMonitorsDescribePOST(c echo.Context) error {
	ctx := context.Background()
	monitorSlug := c.Param("slug")

	update := UpdateMonitor{
		WorkerGroups: strings.TrimSpace(c.FormValue("workergroups")),
		Schedule:     c.FormValue("schedule"),
		Script:       c.FormValue("script"),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(update)
	if err != nil {
		return err
	}

	monitor, err := services.GetMonitor(ctx, h.db, monitorSlug)
	if err != nil {
		return err
	}
	monitor.Schedule = update.Schedule
	monitor.Script = update.Script

	err = services.UpdateMonitor(
		ctx,
		h.db,
		monitor,
	)
	if err != nil {
		return err
	}

	workerGroups := []*models.WorkerGroup{}
	for _, group := range strings.Split(update.WorkerGroups, " ") {
		if group == "" {
			continue
		}
		workerGroup, err := services.GetWorkerGroup(ctx, h.db, slug.Make(group))
		if err != nil {
			if err == sql.ErrNoRows {
				workerGroup = &models.WorkerGroup{Name: group, Slug: slug.Make(group)}
				err = services.CreateWorkerGroup(ctx, h.db, workerGroup)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		workerGroups = append(workerGroups, workerGroup)
	}

	err = services.UpdateMonitorWorkerGroups(ctx, h.db, monitor, workerGroups)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateMonitorSchedule(ctx, h.temporal, monitor, workerGroups)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/monitors/%s", monitorSlug))
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
	monitorSlug := slug.Make(c.FormValue("name"))

	create := CreateMonitor{
		Name:         c.FormValue("name"),
		Schedule:     c.FormValue("schedule"),
		WorkerGroups: c.FormValue("workergroups"),
		Script:       c.FormValue("script"),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(create)
	if err != nil {
		return err
	}

	workerGroups := []*models.WorkerGroup{}
	for _, group := range strings.Split(create.WorkerGroups, " ") {
		if group == "" {
			continue
		}
		workerGroup, err := services.GetWorkerGroup(ctx, h.db, slug.Make(group))
		if err != nil {
			if err == sql.ErrNoRows {
				workerGroup = &models.WorkerGroup{Name: group, Slug: slug.Make(group)}
				err = services.CreateWorkerGroup(ctx, h.db, workerGroup)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		workerGroups = append(workerGroups, workerGroup)
	}

	monitor := &models.Monitor{
		Name:     create.Name,
		Slug:     monitorSlug,
		Schedule: create.Schedule,
		Script:   create.Script,
	}

	err = services.CreateMonitor(
		ctx,
		h.db,
		monitor,
	)
	if err != nil {
		return err
	}

	err = services.UpdateMonitorWorkerGroups(ctx, h.db, monitor, workerGroups)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateMonitorSchedule(ctx, h.temporal, monitor, workerGroups)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/monitors")
}
