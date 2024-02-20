package templates

import (
	"embed"
	"io"
	"log"
	"text/template"

	"github.com/labstack/echo/v4"
)

//go:embed *
var templates embed.FS

const base = "components/base.tmpl"

type Templates struct {
	templates map[string]*template.Template
}

func load(files ...string) *template.Template {
	files = append(files, base)
	return template.Must(template.ParseFS(templates, files...))
}

func loadSettings(files ...string) *template.Template {
	files = append(files, "components/settings.tmpl")
	return load(files...)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: map[string]*template.Template{
			"404.tmpl":                            load("pages/404.tmpl"),
			"index.tmpl":                          load("pages/index.tmpl"),
			"settings_overview.tmpl":              loadSettings("pages/settings_overview.tmpl"),
			"settings_workers.tmpl":               loadSettings("pages/settings_workers.tmpl"),
			"settings_workers_create.tmpl":        loadSettings("pages/settings_workers_create.tmpl"),
			"settings_healthchecks.tmpl":          loadSettings("pages/settings_healthchecks.tmpl"),
			"settings_healthchecks_create.tmpl":   loadSettings("pages/settings_healthchecks_create.tmpl"),
			"settings_healthchecks_describe.tmpl": loadSettings("pages/settings_healthchecks_describe.tmpl"),
		},
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.templates[name] == nil {
		log.Printf("template not found: %s", name)
		return echo.ErrNotFound
	}

	err := t.templates[name].ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("error rendering template: %s", err)
	}

	return err
}
