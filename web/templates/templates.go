package templates

import (
	"embed"
	"io"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mentos1386/zdravko/pkg/script"
)

//go:embed *
var templates embed.FS

const base = "components/base.tmpl"

type Templates struct {
	templates map[string]*template.Template
}

func load(files ...string) *template.Template {
	files = append(files, base)

	t := template.New("default").Funcs(
		template.FuncMap{
			"MathDivide": func(a, b int) float64 {
				if b == 0 {
					return 0
				}
				return float64(a) / float64(b)
			},
			"DurationRoundSecond": func(d time.Duration) time.Duration {
				return d.Round(time.Second)
			},
			"DurationRoundMillisecond": func(d time.Duration) time.Duration {
				return d.Round(time.Millisecond)
			},
			"StringsJoin":          strings.Join,
			"Now":                  time.Now,
			"ScriptUnescapeString": script.UnescapeString,
			"ScriptEscapeString":   script.EscapeString,
		})

	return template.Must(t.ParseFS(templates, files...))
}

func loadSettings(files ...string) *template.Template {
	files = append(files, "components/settings.tmpl")
	return load(files...)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: map[string]*template.Template{
			"404.tmpl":                             load("pages/404.tmpl"),
			"index.tmpl":                           load("pages/index.tmpl"),
			"incidents.tmpl":                       load("pages/incidents.tmpl"),
			"settings_home.tmpl":                   loadSettings("pages/settings_home.tmpl"),
			"settings_triggers.tmpl":               loadSettings("pages/settings_triggers.tmpl"),
			"settings_triggers_create.tmpl":        loadSettings("pages/settings_triggers_create.tmpl"),
			"settings_triggers_describe.tmpl":      loadSettings("pages/settings_triggers_describe.tmpl"),
			"settings_targets.tmpl":                loadSettings("pages/settings_targets.tmpl"),
			"settings_targets_create.tmpl":         loadSettings("pages/settings_targets_create.tmpl"),
			"settings_targets_describe.tmpl":       loadSettings("pages/settings_targets_describe.tmpl"),
			"settings_incidents.tmpl":              loadSettings("pages/settings_incidents.tmpl"),
			"settings_notifications.tmpl":          loadSettings("pages/settings_notifications.tmpl"),
			"settings_worker_groups.tmpl":          loadSettings("pages/settings_worker_groups.tmpl"),
			"settings_worker_groups_create.tmpl":   loadSettings("pages/settings_worker_groups_create.tmpl"),
			"settings_worker_groups_describe.tmpl": loadSettings("pages/settings_worker_groups_describe.tmpl"),
			"settings_checks.tmpl":                 loadSettings("pages/settings_checks.tmpl"),
			"settings_checks_create.tmpl":          loadSettings("pages/settings_checks_create.tmpl"),
			"settings_checks_describe.tmpl":        loadSettings("pages/settings_checks_describe.tmpl"),
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
