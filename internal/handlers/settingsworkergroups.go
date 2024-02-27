package handlers

import (
	"context"
	"fmt"
	"net/http"

	"code.tjo.space/mentos1386/zdravko/database/models"
	"code.tjo.space/mentos1386/zdravko/internal/jwt"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

type WorkerWithTokenAndActiveWorkers struct {
	*models.WorkerGroup
	Token         string
	ActiveWorkers []string
}

type WorkerGroupWithActiveWorkers struct {
	*models.WorkerGroupWithMonitors
	ActiveWorkers []string
}

type SettingsWorkerGroups struct {
	*Settings
	WorkerGroups       []*WorkerGroupWithActiveWorkers
	WorkerGroupsLength int
}

type SettingsWorker struct {
	*Settings
	Worker *WorkerWithTokenAndActiveWorkers
}

func (h *BaseHandler) SettingsWorkerGroupsGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	workerGroups, err := services.GetWorkerGroupsWithMonitors(context.Background(), h.db)
	if err != nil {
		return err
	}

	workerGroupsWithActiveWorkers := make([]*WorkerGroupWithActiveWorkers, len(workerGroups))
	for i, workerGroup := range workerGroups {
		activeWorkers, err := services.GetActiveWorkers(context.Background(), workerGroup.Slug, h.temporal)
		if err != nil {
			return err
		}
		workerGroupsWithActiveWorkers[i] = &WorkerGroupWithActiveWorkers{
			WorkerGroupWithMonitors: workerGroup,
			ActiveWorkers:           activeWorkers,
		}
	}

	return c.Render(http.StatusOK, "settings_worker_groups.tmpl", &SettingsWorkerGroups{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Worker Groups"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Worker Groups")},
		),
		WorkerGroups:       workerGroupsWithActiveWorkers,
		WorkerGroupsLength: len(workerGroupsWithActiveWorkers),
	})
}

func (h *BaseHandler) SettingsWorkerGroupsDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("slug")

	worker, err := services.GetWorkerGroup(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	// Allow write access to default namespace
	token, err := jwt.NewTokenForWorker(h.config.Jwt.PrivateKey, h.config.Jwt.PublicKey, worker)
	if err != nil {
		return err
	}

	activeWorkers, err := services.GetActiveWorkers(context.Background(), worker.Slug, h.temporal)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_worker_groups_describe.tmpl", &SettingsWorker{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Worker Groups"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Worker Groups"),
				{
					Path:       fmt.Sprintf("/settings/worker-groups/%s", slug),
					Title:      "Describe",
					Breadcrumb: worker.Name,
				},
			}),
		Worker: &WorkerWithTokenAndActiveWorkers{
			WorkerGroup:   worker,
			Token:         token,
			ActiveWorkers: activeWorkers,
		},
	})
}

func (h *BaseHandler) SettingsWorkerGroupsCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_worker_groups_create.tmpl", NewSettings(
		cc.Principal.User,
		GetPageByTitle(SettingsPages, "Worker Groups"),
		[]*components.Page{
			GetPageByTitle(SettingsPages, "Worker Groups"),
			GetPageByTitle(SettingsPages, "Worker Groups Create"),
		},
	))
}

func (h *BaseHandler) SettingsWorkerGroupsCreatePOST(c echo.Context) error {
	ctx := context.Background()
	slug := slug.Make(c.FormValue("name"))

	workerGroup := &models.WorkerGroup{
		Name: c.FormValue("name"),
		Slug: slug,
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(workerGroup)
	if err != nil {
		return err
	}

	err = services.CreateWorkerGroup(
		ctx,
		h.db,
		workerGroup,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/worker-groups/%s", slug))
}
