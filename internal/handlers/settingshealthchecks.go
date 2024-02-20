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
	Healthchecks       []*models.HealthcheckHttp
	HealthchecksLength int
}

type SettingsHealthcheck struct {
	*Settings
	Healthcheck *models.HealthcheckHttp
}

func (h *BaseHandler) SettingsHealthchecksGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	healthchecks, err := h.query.HealthcheckHttp.WithContext(context.Background()).Find()
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

	healthcheck, err := services.GetHealthcheckHttp(context.Background(), h.query, slug)
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

	healthcheckHttp := &models.HealthcheckHttp{
		Healthcheck: models.Healthcheck{
			Name:         c.FormValue("name"),
			Slug:         slug.Make(c.FormValue("name")),
			Schedule:     c.FormValue("schedule"),
			WorkerGroups: strings.Split(c.FormValue("workergroups"), ","),
		},
		Url:    c.FormValue("url"),
		Method: c.FormValue("method"),
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(healthcheckHttp)
	if err != nil {
		return err
	}

	err = services.CreateHealthcheckHttp(
		ctx,
		h.db,
		healthcheckHttp,
	)
	if err != nil {
		return err
	}

	err = services.StartHealthcheckHttp(ctx, h.temporal, healthcheckHttp)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/healthchecks")
}
