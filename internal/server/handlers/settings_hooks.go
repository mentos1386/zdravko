package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/pkg/script"
	"github.com/mentos1386/zdravko/web/templates/components"
)

type CreateHook struct {
	Name     string `validate:"required"`
	Schedule string `validate:"required,cron"`
	Script   string `validate:"required"`
}

type UpdateHook struct {
	Schedule string `validate:"required,cron"`
	Script   string `validate:"required"`
}

type HookWithState struct {
	*models.Hook
	State models.HookState
}

type SettingsHooks struct {
	*Settings
	Hooks   []*HookWithState
	History []struct {
		CreatedAt time.Time
		Status    string
		Note      string
	}
}

type SettingsHook struct {
	*Settings
	Hook    *HookWithState
	History []*services.HookHistory
}

type SettingsHookCreate struct {
	*Settings
	ExampleScript string
	ExampleFilter string
}

func (h *BaseHandler) SettingsHooksGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	hooks, err := services.GetHooks(context.Background(), h.db)
	if err != nil {
		return err
	}

	hooksWithState := make([]*HookWithState, len(hooks))
	for i, hook := range hooks {
		state, err := services.GetHookState(context.Background(), h.temporal, hook.Id)
		if err != nil {
			h.logger.Error("Failed to get hook state", "error", err)
			state = models.HookStateUnknown
		}
		hooksWithState[i] = &HookWithState{
			Hook:  hook,
			State: state,
		}
	}

	return c.Render(http.StatusOK, "settings_hooks.tmpl", &SettingsHooks{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Hooks"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Hooks")},
		),
		Hooks: hooksWithState,
	})
}

func (h *BaseHandler) SettingsHooksDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("id")

	hook, err := services.GetHook(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	status, err := services.GetHookState(context.Background(), h.temporal, hook.Id)
	if err != nil {
		return err
	}

	hookWithStatus := &HookWithState{
		Hook:  hook,
		State: status,
	}

	history, err := services.GetHookHistoryForHook(context.Background(), h.temporal, slug)
	if err != nil {
		return err
	}

	maxElements := 10
	if len(history) < maxElements {
		maxElements = len(history)
	}

	return c.Render(http.StatusOK, "settings_hooks_describe.tmpl", &SettingsHook{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Hooks"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Hooks"),
				{
					Path:       fmt.Sprintf("/settings/hooks/%s", slug),
					Title:      "Describe",
					Breadcrumb: hook.Name,
				},
			}),
		Hook:    hookWithStatus,
		History: history[:maxElements],
	})
}

func (h *BaseHandler) SettingsHooksDescribeDELETE(c echo.Context) error {
	slug := c.Param("id")

	err := services.DeleteHook(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.DeleteHookSchedule(context.Background(), h.temporal, slug)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/hooks")
}

func (h *BaseHandler) SettingsHooksDisableGET(c echo.Context) error {
	slug := c.Param("id")

	hook, err := services.GetHook(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.SetHookState(context.Background(), h.temporal, hook.Id, models.HookStatePaused)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/hooks/%s", slug))
}

func (h *BaseHandler) SettingsHooksEnableGET(c echo.Context) error {
	slug := c.Param("id")

	hook, err := services.GetHook(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.SetHookState(context.Background(), h.temporal, hook.Id, models.HookStateActive)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/hooks/%s", slug))
}

func (h *BaseHandler) SettingsHooksDescribePOST(c echo.Context) error {
	ctx := context.Background()
	hookId := c.Param("id")

	update := UpdateHook{
		Schedule: c.FormValue("schedule"),
		Script:   script.EscapeString(c.FormValue("script")),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(update)
	if err != nil {
		return err
	}

	hook, err := services.GetHook(ctx, h.db, hookId)
	if err != nil {
		return err
	}
	hook.Schedule = update.Schedule
	hook.Script = update.Script

	err = services.UpdateHook(
		ctx,
		h.db,
		hook,
	)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateHookSchedule(ctx, h.temporal, hook)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/hooks/%s", hookId))
}

func (h *BaseHandler) SettingsHooksCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_hooks_create.tmpl", &SettingsHookCreate{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Hooks"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Hooks"),
				GetPageByTitle(SettingsPages, "Hooks Create"),
			},
		),
		ExampleScript: h.examples.Hook,
	})
}

func (h *BaseHandler) SettingsHooksCreatePOST(c echo.Context) error {
	ctx := context.Background()
	hookId := slug.Make(c.FormValue("name"))

	create := CreateHook{
		Name:     c.FormValue("name"),
		Schedule: c.FormValue("schedule"),
		Script:   script.EscapeString(c.FormValue("script")),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(create)
	if err != nil {
		return err
	}

	hook := &models.Hook{
		Name:     create.Name,
		Id:       hookId,
		Schedule: create.Schedule,
		Script:   create.Script,
	}

	err = services.CreateHook(
		ctx,
		h.db,
		hook,
	)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateHookSchedule(ctx, h.temporal, hook)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/hooks")
}
