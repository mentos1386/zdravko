package handlers

import (
	"net/http"

	"github.com/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type Notification struct{}

type SettingsNotifications struct {
	*Settings
	Notifications []*Notification
}

func (h *BaseHandler) SettingsNotificationsGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	notifications := make([]*Notification, 0)

	return c.Render(http.StatusOK, "settings_notifications.tmpl", &SettingsNotifications{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Notifications"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Notifications")},
		),
		Notifications: notifications,
	})
}
