package handlers

import (
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/web/templates"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
)

type Settings struct {
	*components.Base
	SettingsPage  *components.Page
	SettingsPages []*components.Page
	User          *AuthenticatedUser
}

var SettingsPages = []*components.Page{
	{Path: "/settings", Title: "Overview"},
	{Path: "/settings/healthchecks", Title: "Healthchecks"},
	{Path: "/settings/workers", Title: "Workers"},
	{Path: "/temporal", Title: "Temporal"},
	{Path: "/oauth2/logout", Title: "Logout"},
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

	err = ts.ExecuteTemplate(w, "base", &Settings{
		Base: &components.Base{
			Page:  GetPageByTitle(Pages, "Settings"),
			Pages: Pages,
		},
		SettingsPage:  GetPageByTitle(SettingsPages, "Overview"),
		SettingsPages: SettingsPages,
		User:          user,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BaseHandler) SettingsHealthchecksGET(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_healthchecks.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", &Settings{
		Base: &components.Base{
			Page:  GetPageByTitle(Pages, "Settings"),
			Pages: Pages,
		},
		SettingsPage:  GetPageByTitle(SettingsPages, "Healthchecks"),
		SettingsPages: SettingsPages,
		User:          user,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
