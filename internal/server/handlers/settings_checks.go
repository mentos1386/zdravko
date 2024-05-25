package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/pkg/script"
	"github.com/mentos1386/zdravko/web/templates/components"
)

type CreateCheck struct {
	Name         string `validate:"required"`
	WorkerGroups string `validate:"required"`
	Schedule     string `validate:"required,cron"`
	Script       string `validate:"required"`
	Filter       string `validate:"required"`
}

type UpdateCheck struct {
	WorkerGroups string `validate:"required"`
	Schedule     string `validate:"required,cron"`
	Script       string `validate:"required"`
	Filter       string `validate:"required"`
}

type CheckWithWorkerGroupsAndState struct {
	*models.CheckWithWorkerGroups
	State models.CheckState
}

type SettingsChecks struct {
	*Settings
	Checks  []*CheckWithWorkerGroupsAndState
	History []struct {
		CreatedAt time.Time
		Status    string
		Note      string
	}
}

type SettingsCheck struct {
	*Settings
	Check   *CheckWithWorkerGroupsAndState
	History []*services.CheckHistory
}

type SettingsCheckCreate struct {
	*Settings
	ExampleScript string
	ExampleFilter string
}

func (h *BaseHandler) SettingsChecksGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	checks, err := services.GetChecksWithWorkerGroups(context.Background(), h.db)
	if err != nil {
		return err
	}

	checksWithState := make([]*CheckWithWorkerGroupsAndState, len(checks))
	for i, check := range checks {
		state, err := services.GetCheckState(context.Background(), h.temporal, check.Id)
		if err != nil {
			h.logger.Error("Failed to get check state", "error", err)
			state = models.CheckStateUnknown
		}
		checksWithState[i] = &CheckWithWorkerGroupsAndState{
			CheckWithWorkerGroups: check,
			State:                 state,
		}
	}

	return c.Render(http.StatusOK, "settings_checks.tmpl", &SettingsChecks{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Checks"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Checks")},
		),
		Checks: checksWithState,
	})
}

func (h *BaseHandler) SettingsChecksDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("id")

	check, err := services.GetCheckWithWorkerGroups(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	status, err := services.GetCheckState(context.Background(), h.temporal, check.Id)
	if err != nil {
		return err
	}

	checkWithStatus := &CheckWithWorkerGroupsAndState{
		CheckWithWorkerGroups: check,
		State:                 status,
	}

	history, err := services.GetCheckHistoryForCheck(context.Background(), h.temporal, slug)
	if err != nil {
		return err
	}

	maxElements := 10
	if len(history) < maxElements {
		maxElements = len(history)
	}

	return c.Render(http.StatusOK, "settings_checks_describe.tmpl", &SettingsCheck{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Checks"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Checks"),
				{
					Path:       fmt.Sprintf("/settings/checks/%s", slug),
					Title:      "Describe",
					Breadcrumb: check.Name,
				},
			}),
		Check:   checkWithStatus,
		History: history[:maxElements],
	})
}

func (h *BaseHandler) SettingsChecksDescribeDELETE(c echo.Context) error {
	slug := c.Param("id")

	err := services.DeleteCheck(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.DeleteCheckSchedule(context.Background(), h.temporal, slug)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/checks")
}

func (h *BaseHandler) SettingsChecksDisableGET(c echo.Context) error {
	slug := c.Param("id")

	check, err := services.GetCheck(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.SetCheckState(context.Background(), h.temporal, check.Id, models.CheckStatePaused)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/checks/%s", slug))
}

func (h *BaseHandler) SettingsChecksEnableGET(c echo.Context) error {
	slug := c.Param("id")

	check, err := services.GetCheck(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.SetCheckState(context.Background(), h.temporal, check.Id, models.CheckStateActive)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/checks/%s", slug))
}

func (h *BaseHandler) SettingsChecksDescribePOST(c echo.Context) error {
	ctx := context.Background()
	checkId := c.Param("id")

	update := UpdateCheck{
		WorkerGroups: strings.ToLower(strings.TrimSpace(c.FormValue("workergroups"))),
		Schedule:     c.FormValue("schedule"),
		Script:       script.EscapeString(c.FormValue("script")),
		Filter:       c.FormValue("filter"),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(update)
	if err != nil {
		return err
	}

	check, err := services.GetCheck(ctx, h.db, checkId)
	if err != nil {
		return err
	}
	check.Schedule = update.Schedule
	check.Script = update.Script
	check.Filter = update.Filter

	err = services.UpdateCheck(
		ctx,
		h.db,
		check,
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
				workerGroup = &models.WorkerGroup{Name: group, Id: slug.Make(group)}
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

	err = services.UpdateCheckWorkerGroups(ctx, h.db, check, workerGroups)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateCheckSchedule(ctx, h.temporal, check, workerGroups)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/checks/%s", checkId))
}

func (h *BaseHandler) SettingsChecksCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_checks_create.tmpl", &SettingsCheckCreate{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Checks"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Checks"),
				GetPageByTitle(SettingsPages, "Checks Create"),
			},
		),
		ExampleScript: h.examples.Check,
		ExampleFilter: h.examples.Filter,
	})
}

func (h *BaseHandler) SettingsChecksCreatePOST(c echo.Context) error {
	ctx := context.Background()
	checkId := slug.Make(c.FormValue("name"))

	create := CreateCheck{
		Name:         c.FormValue("name"),
		WorkerGroups: strings.ToLower(strings.TrimSpace(c.FormValue("workergroups"))),
		Schedule:     c.FormValue("schedule"),
		Script:       script.EscapeString(c.FormValue("script")),
		Filter:       c.FormValue("filter"),
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
				workerGroup = &models.WorkerGroup{Name: group, Id: slug.Make(group)}
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

	check := &models.Check{
		Name:     create.Name,
		Id:       checkId,
		Schedule: create.Schedule,
		Script:   create.Script,
		Filter:   create.Filter,
	}

	err = services.CreateCheck(
		ctx,
		h.db,
		check,
	)
	if err != nil {
		return err
	}

	err = services.UpdateCheckWorkerGroups(ctx, h.db, check, workerGroups)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateCheckSchedule(ctx, h.temporal, check, workerGroups)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/checks")
}
