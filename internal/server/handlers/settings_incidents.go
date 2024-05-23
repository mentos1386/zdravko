package handlers

import (
	"net/http"

	"github.com/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type Incident struct{}

type SettingsIncidents struct {
	*Settings
	Incidents []*Incident
}

func (h *BaseHandler) SettingsIncidentsGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	incidents := make([]*Incident, 0)

	return c.Render(http.StatusOK, "settings_incidents.tmpl", &SettingsIncidents{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Incidents"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Incidents")},
		),
		Incidents: incidents,
	})
}
