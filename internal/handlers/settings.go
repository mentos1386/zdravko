package handlers

import (
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type Settings struct {
	*components.Base
	SettingsSidebarActive *components.Page
	SettingsSidebar       []*components.Page
	User                  *AuthenticatedUser
	SettingsBreadcrumbs   []*components.Page
}

func NewSettings(user *AuthenticatedUser, page *components.Page, breadCrumbs []*components.Page) *Settings {
	return &Settings{
		Base: &components.Base{
			NavbarActive: GetPageByTitle(Pages, "Settings"),
			Navbar:       Pages,
		},
		SettingsSidebarActive: page,
		SettingsSidebar:       SettingsNavbar,
		SettingsBreadcrumbs:   breadCrumbs,
		User:                  user,
	}
}

var SettingsPages = []*components.Page{
	{Path: "/settings", Title: "Overview", Breadcrumb: "Overview"},
	{Path: "/settings/monitors", Title: "Monitors", Breadcrumb: "Monitors"},
	{Path: "/settings/monitors/create", Title: "Monitors Create", Breadcrumb: "Create"},
	{Path: "/settings/worker-groups", Title: "Worker Groups", Breadcrumb: "Worker Groups"},
	{Path: "/settings/worker-groups/create", Title: "Worker Groups Create", Breadcrumb: "Create"},
	{Path: "/settings/notifications", Title: "Notifications", Breadcrumb: "Notifications"},
	{Path: "/settings/notifications/create", Title: "Notifications Create", Breadcrumb: "Create"},
	{Path: "/settings/incidents", Title: "Incidents", Breadcrumb: "Incidents"},
	{Path: "/settings/incidents/create", Title: "Incidents Create", Breadcrumb: "Create"},
	{Path: "/settings/temporal", Title: "Temporal", Breadcrumb: "Temporal"},
	{Path: "/oauth2/logout", Title: "Logout", Breadcrumb: "Logout"},
}

var SettingsNavbar = []*components.Page{
	GetPageByTitle(SettingsPages, "Overview"),
	GetPageByTitle(SettingsPages, "Incidents"),
	GetPageByTitle(SettingsPages, "Monitors"),
	GetPageByTitle(SettingsPages, "Notifications"),
	GetPageByTitle(SettingsPages, "Worker Groups"),
	GetPageByTitle(SettingsPages, "Temporal"),
	GetPageByTitle(SettingsPages, "Logout"),
}

type SettingsOverview struct {
	*Settings
	WorkerGroupsCount  int
	MonitorsCount      int
	NotificationsCount int
	History            []*services.MonitorHistoryWithMonitor
}

func (h *BaseHandler) SettingsOverviewGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)
	ctx := c.Request().Context()

	workerGroups, err := services.CountWorkerGroups(ctx, h.db)
	if err != nil {
		return err
	}

	monitors, err := services.CountMonitors(ctx, h.db)
	if err != nil {
		return err
	}

	history, err := services.GetLastNMonitorHistory(ctx, h.db, 10)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "settings_overview.tmpl", SettingsOverview{
		Settings: NewSettings(
			cc.Principal.User,
			GetPageByTitle(SettingsPages, "Overview"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Overview")},
		),
		WorkerGroupsCount:  workerGroups,
		MonitorsCount:      monitors,
		NotificationsCount: 42,
		History:            history,
	})
}
