package handlers

import (
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/web/templates"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
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
	{Path: "/settings/healthchecks", Title: "Healthchecks", Breadcrumb: "Healthchecks"},
	{Path: "/settings/healthchecks/create", Title: "Healthchecks Create", Breadcrumb: "Create"},
	{Path: "/settings/cronjobs", Title: "Cronjobs", Breadcrumb: "Cronjobs"},
	{Path: "/settings/workers", Title: "Workers", Breadcrumb: "Workers"},
	{Path: "/settings/workers/create", Title: "Workers Create", Breadcrumb: "Create"},
	{Path: "/temporal", Title: "Temporal", Breadcrumb: "Temporal"},
	{Path: "/oauth2/logout", Title: "Logout", Breadcrumb: "Logout"},
}

var SettingsNavbar = []*components.Page{
	GetPageByTitle(SettingsPages, "Overview"),
	GetPageByTitle(SettingsPages, "Healthchecks"),
	GetPageByTitle(SettingsPages, "Cronjobs"),
	GetPageByTitle(SettingsPages, "Workers"),
	GetPageByTitle(SettingsPages, "Temporal"),
	GetPageByTitle(SettingsPages, "Logout"),
}

func (h *BaseHandler) SettingsOverviewGET(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_overview.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", NewSettings(
		user,
		GetPageByTitle(SettingsPages, "Overview"),
		[]*components.Page{GetPageByTitle(SettingsPages, "Overview")},
	))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
