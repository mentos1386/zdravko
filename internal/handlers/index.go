package handlers

import (
	"math/rand"
	"net/http"
	"text/template"

	"code.tjo.space/mentos1386/zdravko/web/templates"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
)

type IndexData struct {
	*components.Base
	HealthChecks []*HealthCheck
}

type HealthCheck struct {
	Domain  string
	Healthy bool
	Uptime  string
	History []bool
}

func newMockHealthCheck(domain string) *HealthCheck {
	randBool := func() bool {
		return rand.Intn(2) == 1
	}

	var history []bool
	for i := 0; i < 90; i++ {
		history = append(history, randBool())
	}

	return &HealthCheck{
		Domain:  domain,
		Healthy: randBool(),
		Uptime:  "100",
		History: history,
	}
}

func (h *BaseHandler) Index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(templates.Templates,
		"components/base.tmpl",
		"pages/index.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", &IndexData{
		Base: &components.Base{
			Page:  GetPageByTitle(Pages, "Status"),
			Pages: Pages,
		},
		HealthChecks: []*HealthCheck{
			newMockHealthCheck("example.com"),
			newMockHealthCheck("example.org"),
			newMockHealthCheck("example.net"),
			newMockHealthCheck("foo.example.net"),
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
