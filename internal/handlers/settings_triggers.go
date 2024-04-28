package handlers

import (
	"context"
	"fmt"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

type Incident struct{}

type SettingsIncidents struct {
	*Settings
	Incidents []*Incident
}

// FIXME: This should be moved in to a files to not worry about so much escaping.
const exampleTrigger = `
import kv from 'zdravko/kv';
import slack from 'zdravko/notify/slack';

export default function (monitor, outcome) {
  // If the outcome is not failure, we can reset the counter.
  if (outcome.status !== 'FAILURE') {
    return kv.delete(` + "\\`\\${monitor.name}:issues:5min\\`" + `);
  }

  const count = kv.get(` + "\\`\\${monitor.name}:issues:5min\\`" + `) || 0;

  if (count > 5) {
    slack.notify(` + "\\`\\${monitor.name} has had more than 5 issues in the last 5 minutes\\`" + `);
  }

  // Increment and set TTL to 5 minutes
  kv.increment(` + "\\`\\${monitor.name}:issues:5min\\`" + `, count + 1);
}
`

type CreateTrigger struct {
	Name   string `validate:"required"`
	Script string `validate:"required"`
}

type UpdateTrigger struct {
	Script string `validate:"required"`
}

type SettingsTriggers struct {
	*Settings
	Triggers []*models.Trigger
}

type SettingsTrigger struct {
	*Settings
	Trigger *models.Trigger
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
	return c.Render(http.StatusOK, "settings_triggers.tmpl", &SettingsTriggers{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Triggers"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Triggers")},
		),
		Triggers: triggers,
	})
}

func (h *BaseHandler) SettingsTriggersDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("id")

	trigger, err := services.GetTrigger(context.Background(), h.db, slug)
	if err != nil {
		return err
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
		Trigger: trigger,
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

	trigger.Status = services.TriggerStatusPaused
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

	trigger.Status = services.TriggerStatusActive
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
		Script: c.FormValue("script"),
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
		Example: exampleTrigger,
	})
}

func (h *BaseHandler) SettingsTriggersCreatePOST(c echo.Context) error {
	ctx := context.Background()
	triggerId := slug.Make(c.FormValue("name"))

	create := CreateTrigger{
		Name:   c.FormValue("name"),
		Script: c.FormValue("script"),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(create)
	if err != nil {
		return err
	}

	trigger := &models.Trigger{
		Name:   create.Name,
		Id:     triggerId,
		Script: create.Script,
		Status: services.TriggerStatusActive,
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
