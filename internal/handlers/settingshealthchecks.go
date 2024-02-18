package handlers

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

type SettingsHealthchecks struct {
	*Settings
	Healthchecks       []*models.HealthcheckHttp
	HealthchecksLength int
}

type SettingsHealthcheck struct {
	*Settings
	Healthcheck *models.HealthcheckHttp
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

	healthchecks, err := h.query.HealthcheckHttp.WithContext(context.Background()).Find()
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
	slug := vars["slug"]

	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_healthchecks_describe.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	healthcheck, err := services.GetHealthcheckHttp(context.Background(), h.query, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", &SettingsHealthcheck{
		Settings: NewSettings(
			user,
			GetPageByTitle(SettingsPages, "Healthchecks"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Healthchecks"),
				{
					Path:       fmt.Sprintf("/settings/healthchecks/%s", slug),
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
	ctx := context.Background()

	healthcheckHttp := &models.HealthcheckHttp{
		Healthcheck: models.Healthcheck{
			Name:     r.FormValue("name"),
			Slug:     slug.Make(r.FormValue("name")),
			Schedule: r.FormValue("schedule"),
		},
		Url:    r.FormValue("url"),
		Method: r.FormValue("method"),
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(healthcheckHttp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = services.CreateHealthcheckHttp(
		ctx,
		h.db,
		healthcheckHttp,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = services.StartHealthcheckHttp(ctx, h.temporal, healthcheckHttp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/settings/healthchecks", http.StatusSeeOther)
}
