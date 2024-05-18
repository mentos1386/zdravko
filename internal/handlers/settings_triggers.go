package handlers

import (
	"context"
	"fmt"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/script"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

type CreateTrigger struct {
	Name   string `validate:"required"`
	Script string `validate:"required"`
}

type UpdateTrigger struct {
	Script string `validate:"required"`
}

type TriggerWithState struct {
	*models.Trigger
	State models.TriggerState
}

type SettingsTriggers struct {
	*Settings
	Triggers []*TriggerWithState
}

type SettingsTrigger struct {
	*Settings
	Trigger *TriggerWithState
	History []*models.TriggerHistory
}

type SettingsTriggerCreate struct {
	*Settings
	Example string
}

func (h *BaseHandler) SettingsTriggersGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	triggers, err := services.GetTriggers(context.Background(), h.db)
	if err != nil {
		return err
	}

	triggersWithState := make([]*TriggerWithState, 0, len(triggers))
	for _, trigger := range triggers {
		triggersWithState = append(triggersWithState, &TriggerWithState{
			Trigger: trigger,
			State:   models.TriggerStateActive,
		})
	}

	return c.Render(http.StatusOK, "settings_triggers.tmpl", &SettingsTriggers{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Triggers"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Triggers")},
		),
		Triggers: triggersWithState,
	})
}

func (h *BaseHandler) SettingsTriggersDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("id")

	trigger, err := services.GetTrigger(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	triggerWithState := &TriggerWithState{
		Trigger: trigger,
		State:   models.TriggerStateActive,
	}

	return c.Render(http.StatusOK, "settings_triggers_describe.tmpl", &SettingsTrigger{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Triggers"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Triggers"),
				{
					Path:       fmt.Sprintf("/settings/triggers/%s", slug),
					Title:      "Describe",
					Breadcrumb: trigger.Name,
				},
			}),
		Trigger: triggerWithState,
	})
}

func (h *BaseHandler) SettingsTriggersDescribeDELETE(c echo.Context) error {
	slug := c.Param("id")

	err := services.DeleteTrigger(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/triggers")
}

func (h *BaseHandler) SettingsTriggersDisableGET(c echo.Context) error {
	slug := c.Param("id")

	trigger, err := services.GetTrigger(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.UpdateTrigger(context.Background(), h.db, trigger)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/triggers/%s", slug))
}

func (h *BaseHandler) SettingsTriggersEnableGET(c echo.Context) error {
	slug := c.Param("id")

	trigger, err := services.GetTrigger(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.UpdateTrigger(context.Background(), h.db, trigger)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/triggers/%s", slug))
}

func (h *BaseHandler) SettingsTriggersDescribePOST(c echo.Context) error {
	ctx := context.Background()
	triggerId := c.Param("id")

	update := UpdateTrigger{
		Script: script.EscapeString(c.FormValue("script")),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(update)
	if err != nil {
		return err
	}

	trigger, err := services.GetTrigger(ctx, h.db, triggerId)
	if err != nil {
		return err
	}
	trigger.Script = update.Script

	err = services.UpdateTrigger(
		ctx,
		h.db,
		trigger,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/triggers/%s", triggerId))
}

func (h *BaseHandler) SettingsTriggersCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_triggers_create.tmpl", &SettingsTriggerCreate{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Triggers"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Triggers"),
				GetPageByTitle(SettingsPages, "Triggers Create"),
			},
		),
		Example: h.examples.Trigger,
	})
}

func (h *BaseHandler) SettingsTriggersCreatePOST(c echo.Context) error {
	ctx := context.Background()
	triggerId := slug.Make(c.FormValue("name"))

	create := CreateTrigger{
		Name:   c.FormValue("name"),
		Script: script.EscapeString(c.FormValue("script")),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(create)
	if err != nil {
		return err
	}

	trigger := &models.Trigger{
		Name:   create.Name,
		Id:     triggerId,
		Script: create.Script,
	}

	err = services.CreateTrigger(
		ctx,
		h.db,
		trigger,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/triggers")
}
