package handlers

import (
	"net/http"

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
	{Path: "/settings/workers", Title: "Workers", Breadcrumb: "Workers"},
	{Path: "/settings/workers/create", Title: "Workers Create", Breadcrumb: "Create"},
	{Path: "/settings/notifications", Title: "Notifications", Breadcrumb: "Notifications"},
	{Path: "/settings/temporal", Title: "Temporal", Breadcrumb: "Temporal"},
	{Path: "/oauth2/logout", Title: "Logout", Breadcrumb: "Logout"},
}

var SettingsNavbar = []*components.Page{
	GetPageByTitle(SettingsPages, "Overview"),
	GetPageByTitle(SettingsPages, "Monitors"),
	GetPageByTitle(SettingsPages, "Workers"),
	GetPageByTitle(SettingsPages, "Notifications"),
	GetPageByTitle(SettingsPages, "Temporal"),
	GetPageByTitle(SettingsPages, "Logout"),
}

func (h *BaseHandler) SettingsOverviewGET(c echo.Context) error {
	cc := c.(AuthenticatedContext)

	return c.Render(http.StatusOK, "settings_overview.tmpl", NewSettings(
		cc.Principal.User,
		GetPageByTitle(SettingsPages, "Overview"),
		[]*components.Page{GetPageByTitle(SettingsPages, "Overview")},
	))
}
