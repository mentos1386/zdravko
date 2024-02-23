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

type SettingsHealthchecks struct {
	*Settings
	Healthchecks       []*models.Healthcheck
	HealthchecksLength int
}

type SettingsHealthcheck struct {
	*Settings
	Healthcheck *models.Healthcheck
}

func (h *BaseHandler) SettingsHealthchecksGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	healthchecks, err := h.query.Healthcheck.WithContext(context.Background()).Find()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_healthchecks.tmpl", &SettingsHealthchecks{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Healthchecks"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Healthchecks")},
		),
		Healthchecks:       healthchecks,
		HealthchecksLength: len(healthchecks),
	})
}

func (h *BaseHandler) SettingsHealthchecksDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("slug")

	healthcheck, err := services.GetHealthcheck(context.Background(), h.query, slug)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_healthchecks_describe.tmpl", &SettingsHealthcheck{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Healthchecks"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Healthchecks"),
				{
					Path:       fmt.Sprintf("/settings/healthchecks/%s", slug),
					Title:      "Describe",
					Breadcrumb: healthcheck.Name,
				},
			}),
		Healthcheck: healthcheck,
	})
}

func (h *BaseHandler) SettingsHealthchecksDescribePOST(c echo.Context) error {
	ctx := context.Background()

	slug := c.Param("slug")

	healthcheck, err := services.GetHealthcheck(ctx, h.query, slug)
	if err != nil {
		return err
	}

	update := &models.Healthcheck{
		Slug:         healthcheck.Slug,
		Name:         healthcheck.Name,
		Schedule:     c.FormValue("schedule"),
		WorkerGroups: strings.Split(c.FormValue("workergroups"), " "),
		Script:       c.FormValue("script"),
	}

	err = validator.New(validator.WithRequiredStructEnabled()).Struct(update)
	if err != nil {
		return err
	}

	err = services.UpdateHealthcheck(
		ctx,
		h.query,
		update,
	)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateHealthcheckSchedule(ctx, h.temporal, healthcheck)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/healthchecks/%s", slug))
}

func (h *BaseHandler) SettingsHealthchecksCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_healthchecks_create.tmpl", NewSettings(
		cc.Principal.User,
		GetPageByTitle(SettingsPages, "Healthchecks"),
		[]*components.Page{
			GetPageByTitle(SettingsPages, "Healthchecks"),
			GetPageByTitle(SettingsPages, "Healthchecks Create"),
		},
	))
}

func (h *BaseHandler) SettingsHealthchecksCreatePOST(c echo.Context) error {
	ctx := context.Background()

	healthcheckHttp := &models.Healthcheck{
		Name:         c.FormValue("name"),
		Slug:         slug.Make(c.FormValue("name")),
		Schedule:     c.FormValue("schedule"),
		WorkerGroups: strings.Split(c.FormValue("workergroups"), " "),
		Script:       c.FormValue("script"),
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(healthcheckHttp)
	if err != nil {
		return err
	}

	err = services.CreateHealthcheck(
		ctx,
		h.query,
		healthcheckHttp,
	)
	if err != nil {
		return err
	}

	err = services.CreateOrUpdateHealthcheckSchedule(ctx, h.temporal, healthcheckHttp)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/healthchecks")
}
