package pages

import (
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/internal/ui"
)

func (p *PageHandler) Settings(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(ui.Templates,
		"components/base.tmpl",
		"pages/settings.tmpl",
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
