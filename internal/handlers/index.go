package handlers

import (
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/web/templates"
)

func (h *BaseHandler) Index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"pages/index.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
