package handlers

import (
	"net/http"

	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/labstack/echo/v4"
)

type SettingsSidebarGroup struct {
	Group string
	Pages []*components.Page
}

type Settings struct {
	*components.Base
	SettingsSidebarActive *components.Page
	SettingsSidebar       []SettingsSidebarGroup
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
		SettingsSidebar:       SettingsSidebar,
		SettingsBreadcrumbs:   breadCrumbs,
		User:                  user,
	}
}

var SettingsPages = []*components.Page{
	{Path: "/settings", Title: "Overview", Breadcrumb: "Overview"},
	{Path: "/settings/incidents", Title: "Incidents", Breadcrumb: "Incidents"},
	{Path: "/settings/targets", Title: "Targets", Breadcrumb: "Targets"},
	{Path: "/settings/targets/create", Title: "Targets Create", Breadcrumb: "Create"},
	{Path: "/settings/hooks", Title: "Hooks", Breadcrumb: "Hooks"},
	{Path: "/settings/hooks/create", Title: "Hooks Create", Breadcrumb: "Create"},
	{Path: "/settings/checks", Title: "Checks", Breadcrumb: "Checks"},
	{Path: "/settings/checks/create", Title: "Checks Create", Breadcrumb: "Create"},
	{Path: "/settings/worker-groups", Title: "Worker Groups", Breadcrumb: "Worker Groups"},
	{Path: "/settings/worker-groups/create", Title: "Worker Groups Create", Breadcrumb: "Create"},
	{Path: "/settings/notifications", Title: "Notifications", Breadcrumb: "Notifications"},
	{Path: "/settings/notifications/create", Title: "Notifications Create", Breadcrumb: "Create"},
	{Path: "/settings/triggers", Title: "Triggers", Breadcrumb: "Triggers"},
	{Path: "/settings/triggers/create", Title: "Triggers Create", Breadcrumb: "Create"},
	{Path: "/settings/temporal", Title: "Temporal", Breadcrumb: "Temporal"},
	{Path: "/oauth2/logout", Title: "Logout", Breadcrumb: "Logout"},
}

var SettingsSidebar = []SettingsSidebarGroup{
	{
		Group: "Overview",
		Pages: []*components.Page{
			GetPageByTitle(SettingsPages, "Overview"),
		},
	},
	{
		Group: "Monitor",
		Pages: []*components.Page{
			GetPageByTitle(SettingsPages, "Targets"),
			GetPageByTitle(SettingsPages, "Checks"),
			GetPageByTitle(SettingsPages, "Hooks"),
		},
	},
	{
		Group: "Alert",
		Pages: []*components.Page{
			GetPageByTitle(SettingsPages, "Triggers"),
		},
	},
	{
		Group: "Notify",
		Pages: []*components.Page{
			GetPageByTitle(SettingsPages, "Incidents"),
			GetPageByTitle(SettingsPages, "Notifications"),
		},
	},
	{
		Group: "Manage",
		Pages: []*components.Page{
			GetPageByTitle(SettingsPages, "Worker Groups"),
			GetPageByTitle(SettingsPages, "Temporal"),
			GetPageByTitle(SettingsPages, "Logout"),
		},
	},
}

type SettingsOverview struct {
	*Settings
	WorkerGroupsCount  int
	ChecksCount        int
	NotificationsCount int
	History            []*services.CheckHistoryWithCheck
}

func (h *BaseHandler) SettingsOverviewGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)
	ctx := c.Request().Context()

	workerGroups, err := services.CountWorkerGroups(ctx, h.db)
	if err != nil {
		return err
	}

	checks, err := services.CountChecks(ctx, h.db)
	if err != nil {
		return err
	}

	history, err := services.GetLastNCheckHistory(ctx, h.db, 10)
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
		ChecksCount:        checks,
		NotificationsCount: 42,
		History:            history,
	})
}
