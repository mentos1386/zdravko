package handlers

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/internal/jwt"
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/services"
	"code.tjo.space/mentos1386/zdravko/web/templates"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

type SettingsWorkers struct {
	*Settings
	Workers       []*models.Worker
	WorkersLength int
}

type SettingsWorker struct {
	*Settings
	Worker *models.Worker
}

func (h *BaseHandler) SettingsWorkersGET(w http.ResponseWriter, r *http.Request, principal *AuthenticatedPrincipal) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_workers.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	workers, err := h.query.Worker.WithContext(context.Background()).Find()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", &SettingsWorkers{
		Settings: NewSettings(
			principal.User,
			GetPageByTitle(SettingsPages, "Workers"),
			[]*components.Page{GetPageByTitle(SettingsPages, "Workers")},
		),
		Workers:       workers,
		WorkersLength: len(workers),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BaseHandler) SettingsWorkersDescribeGET(w http.ResponseWriter, r *http.Request, principal *AuthenticatedPrincipal) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_workers_describe.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	worker, err := services.GetWorker(context.Background(), h.query, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", &SettingsWorker{
		Settings: NewSettings(
			principal.User,
			GetPageByTitle(SettingsPages, "Workers"),
			[]*components.Page{
				GetPageByTitle(SettingsPages, "Workers"),
				{
					Path:       fmt.Sprintf("/settings/workers/%s", slug),
					Title:      "Describe",
					Breadcrumb: worker.Name,
				},
			}),
		Worker: worker,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BaseHandler) SettingsWorkersCreateGET(w http.ResponseWriter, r *http.Request, principal *AuthenticatedPrincipal) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"components/settings.tmpl",
		"pages/settings_workers_create.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", NewSettings(
		principal.User,
		GetPageByTitle(SettingsPages, "Workers"),
		[]*components.Page{
			GetPageByTitle(SettingsPages, "Workers"),
			GetPageByTitle(SettingsPages, "Workers Create"),
		},
	))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BaseHandler) SettingsWorkersCreatePOST(w http.ResponseWriter, r *http.Request, principal *AuthenticatedPrincipal) {
	ctx := context.Background()

	worker := &models.Worker{
		Name:  r.FormValue("name"),
		Slug:  slug.Make(r.FormValue("name")),
		Group: r.FormValue("group"),
	}

	err := validator.New(validator.WithRequiredStructEnabled()).Struct(worker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = services.CreateWorker(
		ctx,
		h.db,
		worker,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/settings/workers", http.StatusSeeOther)
}

func (h *BaseHandler) SettingsWorkersTokenGET(w http.ResponseWriter, r *http.Request, principal *AuthenticatedPrincipal) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	worker, err := services.GetWorker(context.Background(), h.query, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Allow write access to default namespace
	token, err := jwt.NewTokenForWorker(h.config.Jwt.PrivateKey, h.config.Jwt.PublicKey, worker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"token": "` + token + `"}`))
}
