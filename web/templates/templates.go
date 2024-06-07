package templates

import (
	"embed"
	"io"
	"io/fs"
	"log/slog"
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
	logger    *slog.Logger
	templates map[string]*template.Template
}

func load(version string, files ...string) *template.Template {
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
			"Version": func() string {
				return version
			},
		})

	return template.Must(t.ParseFS(templates, files...))
}

func loadSettings(version string, files ...string) *template.Template {
	files = append(files, "components/settings.tmpl")
	return load(version, files...)
}

func NewTemplates(version string, logger *slog.Logger) (*Templates, error) {
	t := Templates{
		logger:    logger,
		templates: map[string]*template.Template{},
	}

	err := fs.WalkDir(templates, "pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, ".tmpl") {
			t.logger.Debug("Loading template", "path", path)
			pathWithoutPrefix := strings.TrimPrefix(path, "pages/")

			if strings.Contains(path, "settings") {
				t.templates[pathWithoutPrefix] = loadSettings(version, path)
			} else {
				t.templates[pathWithoutPrefix] = load(version, path)
			}
		}
		return nil
	})

	return &t, err
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.templates[name] == nil {
		t.logger.Error("template not found", "template", name)
		return echo.ErrNotFound
	}

	err := t.templates[name].ExecuteTemplate(w, "base", data)
	if err != nil {
		t.logger.Error("error rendering template", "template", err)
	}

	return err
}
