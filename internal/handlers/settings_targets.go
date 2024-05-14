package handlers

import (
	"net/http"

	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type Target struct{}

type SettingsTargets struct {
	*Settings
	Targets []*Target
}

func (h *BaseHandler) SettingsTargetsGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	targets := make([]*Target, 0)

	return c.Render(http.StatusOK, "settings_targets.tmpl", &SettingsTargets{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Targets"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Targets")},
		),
		Targets: targets,
	})
}
