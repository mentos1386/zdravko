package handlers

import (
	"context"
	"fmt"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal/jwt"
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

type WorkerWithToken struct {
	*models.Worker
	Token string
}

type SettingsWorkers struct {
	*Settings
	Workers       []*models.Worker
	WorkersLength int
}

type SettingsWorker struct {
	*Settings
	Worker *WorkerWithToken
}

func (h *BaseHandler) SettingsWorkersGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	workers, err := h.query.Worker.WithContext(context.Background()).Find()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_workers.tmpl", &SettingsWorkers{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Workers"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Workers")},
		),
		Workers:       workers,
		WorkersLength: len(workers),
	})
}

func (h *BaseHandler) SettingsWorkersDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("slug")

	worker, err := services.GetWorker(context.Background(), h.query, slug)
	if err != nil {
		return err
	}

	// Allow write access to default namespace
	token, err := jwt.NewTokenForWorker(h.config.Jwt.PrivateKey, h.config.Jwt.PublicKey, worker)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_workers_describe.tmpl", &SettingsWorker{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Workers"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Workers"),
				{
					Path:       fmt.Sprintf("/settings/workers/%s", slug),
					Title:      "Describe",
					Breadcrumb: worker.Name,
				},
			}),
		Worker: &WorkerWithToken{
			Worker: worker,
			Token:  token,
		},
	})
}

func (h *BaseHandler) SettingsWorkersCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_workers_create.tmpl", NewSettings(
		cc.Principal.User,
		GetPageByTitle(SettingsPages, "Workers"),
		[]*components.Page{
			GetPageByTitle(SettingsPages, "Workers"),
			GetPageByTitle(SettingsPages, "Workers Create"),
		},
	))
}

func (h *BaseHandler) SettingsWorkersCreatePOST(c echo.Context) error {
	ctx := context.Background()

	worker := &models.Worker{
		Name:  c.FormValue("name"),
		Slug:  slug.Make(c.FormValue("name")),
		Group: c.FormValue("group"),
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(worker)
	if err != nil {
		return err
	}

	err = services.CreateWorker(
		ctx,
		h.db,
		worker,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/workers")
}
