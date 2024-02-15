package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/web/templates"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/gorilla/mux"
)

type SettingsHealthchecks struct {
	*Settings
	Healthchecks       []*models.HealthcheckHTTP
	HealthchecksLength int
}

type SettingsHealthcheck struct {
	*Settings
	Healthcheck *models.HealthcheckHTTP
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

	healthchecks, err := h.query.HealthcheckHTTP.WithContext(context.Background()).Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", &SettingsHealthchecks{
		Settings: NewSettings(
			user,
			GetPageByTitle(SettingsPages, "Healthchecks"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Healthchecks")},
		),
		Healthchecks:       healthchecks,
		HealthchecksLength: len(healthchecks),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BaseHandler) SettingsHealthchecksDescribeGET(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_healthchecks_describe.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	healthcheck, err := h.query.HealthcheckHTTP.WithContext(context.Background()).Where(
		h.query.HealthcheckHTTP.ID.Eq(uint(id)),
	).First()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", &SettingsHealthcheck{
		Settings: NewSettings(
			user,
			GetPageByTitle(SettingsPages, "Healthchecks"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Healthchecks"),
				&components.Page{
					Path:       fmt.Sprintf("/settings/healthchecks/%d", id),
					Title:      "Describe",
					Breadcrumb: healthcheck.Name,
				},
			}),
		Healthcheck: healthcheck,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BaseHandler) SettingsHealthchecksCreateGET(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_healthchecks_create.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", NewSettings(
		user,
		GetPageByTitle(SettingsPages, "Healthchecks"),
		[]*components.Page{
			GetPageByTitle(SettingsPages, "Healthchecks"),
			GetPageByTitle(SettingsPages, "Healthchecks Create"),
		},
	))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BaseHandler) SettingsHealthchecksCreatePOST(w http.ResponseWriter, r *http.Request, user *AuthenticatedUser) {
	healthcheck := &models.HealthcheckHTTP{
		Healthcheck: models.Healthcheck{
			Name:     r.FormValue("name"),
			Schedule: r.FormValue("schedule"),
		},
		URL:    r.FormValue("url"),
		Method: r.FormValue("method"),
	}
	h.db.Create(healthcheck)

	http.Redirect(w, r, "/settings/healthchecks", http.StatusSeeOther)
}
