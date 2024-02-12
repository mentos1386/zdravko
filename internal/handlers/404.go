package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/web/templates"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
)

func (h *BaseHandler) Error404(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"pages/404.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNotFound)

	err = ts.ExecuteTemplate(w, "base", &components.Base{
		Page:  nil,
		Pages: Pages,
	})
	if err != nil {
		fmt.Println("Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
