package handlers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/database/models"
	"github.com/mentos1386/zdravko/internal/server/services"
	"github.com/mentos1386/zdravko/web/templates/components"
)

type CreateTarget struct {
	Name       string `validate:"required"`
	Group      string `validate:"required"`
	Visibility string `validate:"required,oneof=PUBLIC PRIVATE"`
	Metadata   string `validate:"required"`
}

type UpdateTarget struct {
	Group      string `validate:"required"`
	Visibility string `validate:"required,oneof=PUBLIC PRIVATE"`
	Metadata   string `validate:"required"`
}

type SettingsTargets struct {
	*Settings
	Targets      map[string][]*models.Target
	TargetGroups []string
}

type SettingsTarget struct {
	*Settings
	Target  *models.Target
	History []*services.TargetHistory
}

type SettingsTargetCreate struct {
	*Settings
	Example string
}

func (h *BaseHandler) SettingsTargetsGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	targets, err := services.GetTargets(context.Background(), h.db)
	if err != nil {
		return err
	}

	targetGroups := []string{}
	targetsByGroup := map[string][]*models.Target{}
	for _, target := range targets {
		targetsByGroup[target.Group] = append(targetsByGroup[target.Group], target)
		if !slices.Contains(targetGroups, target.Group) {
			targetGroups = append(targetGroups, target.Group)
		}
	}

	return c.Render(http.StatusOK, "settings_targets.tmpl", &SettingsTargets{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Targets"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Targets")},
		),
		Targets:      targetsByGroup,
		TargetGroups: targetGroups,
	})
}

func (h *BaseHandler) SettingsTargetsDescribeGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	slug := c.Param("id")

	target, err := services.GetTarget(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	history, err := services.GetTargetHistoryForTarget(context.Background(), h.db, slug, services.TargetHistoryDateRange90Minutes)
	if err != nil {
		return err
	}

	maxElements := 10
	if len(history) < maxElements {
		maxElements = len(history)
	}

	return c.Render(http.StatusOK, "settings_targets_describe.tmpl", &SettingsTarget{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Targets"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Targets"),
				{
					Path:       fmt.Sprintf("/settings/targets/%s", slug),
					Title:      "Describe",
					Breadcrumb: target.Name,
				},
			}),
		Target:  target,
		History: history[:maxElements],
	})
}

func (h *BaseHandler) SettingsTargetsDescribeDELETE(c echo.Context) error {
	slug := c.Param("id")

	err := services.DeleteTarget(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/targets")
}

func (h *BaseHandler) SettingsTargetsDisableGET(c echo.Context) error {
	slug := c.Param("id")

	target, err := services.GetTarget(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.SetTargetState(context.Background(), h.db, target.Id, models.TargetStatePaused)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/targets/%s", slug))
}

func (h *BaseHandler) SettingsTargetsEnableGET(c echo.Context) error {
	slug := c.Param("id")

	target, err := services.GetTarget(context.Background(), h.db, slug)
	if err != nil {
		return err
	}

	err = services.SetTargetState(context.Background(), h.db, target.Id, models.TargetStateActive)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/targets/%s", slug))
}

func (h *BaseHandler) SettingsTargetsDescribePOST(c echo.Context) error {
	ctx := context.Background()
	targetId := c.Param("id")

	update := UpdateTarget{
		Group:      strings.ToLower(c.FormValue("group")),
		Visibility: c.FormValue("visibility"),
		Metadata:   c.FormValue("metadata"),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(update)
	if err != nil {
		return err
	}

	target, err := services.GetTarget(ctx, h.db, targetId)
	if err != nil {
		return err
	}
	target.Group = update.Group
	target.Visibility = models.TargetVisibility(update.Visibility)
	target.Metadata = update.Metadata

	err = services.UpdateTarget(
		ctx,
		h.db,
		target,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/settings/targets/%s", targetId))
}

func (h *BaseHandler) SettingsTargetsCreateGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_targets_create.tmpl", &SettingsTargetCreate{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Targets"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Targets"),
				GetPageByTitle(SettingsPages, "Targets Create"),
			},
		),
		Example: h.examples.Target,
	})
}

func (h *BaseHandler) SettingsTargetsCreatePOST(c echo.Context) error {
	ctx := context.Background()
	targetId := slug.Make(c.FormValue("name"))

	create := CreateTarget{
		Name:       c.FormValue("name"),
		Group:      strings.ToLower(c.FormValue("group")),
		Visibility: c.FormValue("visibility"),
		Metadata:   c.FormValue("metadata"),
	}
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(create)
	if err != nil {
		return err
	}

	target := &models.Target{
		Name:       create.Name,
		Group:      create.Group,
		Id:         targetId,
		Visibility: models.TargetVisibility(create.Visibility),
		State:      models.TargetStateActive,
		Metadata:   create.Metadata,
	}

	err = services.CreateTarget(
		ctx,
		h.db,
		target,
	)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/settings/targets")
}
