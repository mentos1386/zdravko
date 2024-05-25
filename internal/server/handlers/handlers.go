package handlers

import (
	"embed"
	"log/slog"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/mentos1386/zdravko/database"
	"github.com/mentos1386/zdravko/internal/config"
	"github.com/mentos1386/zdravko/pkg/script"
	"github.com/mentos1386/zdravko/web/templates/components"
	"go.temporal.io/sdk/client"
	"gopkg.in/yaml.v2"
)

//go:embed examples.yaml
var examplesYaml embed.FS

type examples struct {
	Check   string `yaml:"check"`
	Filter  string `yaml:"filter"`
	Trigger string `yaml:"trigger"`
	Target  string `yaml:"target"`
}

var Pages = []*components.Page{
	{Path: "/", Title: "Status", Breadcrumb: "Status"},
	{Path: "/incidents", Title: "Incidents", Breadcrumb: "Incidents"},
	{Path: "/settings", Title: "Settings", Breadcrumb: "Settings"},
}

func GetPageByTitle(pages []*components.Page, title string) *components.Page {
	for _, p := range pages {
		if p.Title == title {
			return p
		}
	}
	return nil
}

type BaseHandler struct {
	db      *sqlx.DB
	kvStore database.KeyValueStore
	config  *config.ServerConfig
	logger  *slog.Logger

	temporal client.Client

	store *sessions.CookieStore

	examples examples
}

func NewBaseHandler(db *sqlx.DB, kvStore database.KeyValueStore, temporal client.Client, config *config.ServerConfig, logger *slog.Logger) *BaseHandler {
	store := sessions.NewCookieStore([]byte(config.SessionSecret))

	examples := examples{}
	yamlFile, err := examplesYaml.ReadFile("examples.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &examples)
	if err != nil {
		panic(err)
	}

	examples.Check = script.EscapeString(examples.Check)
	examples.Filter = script.EscapeString(examples.Filter)
	examples.Trigger = script.EscapeString(examples.Trigger)
	examples.Target = script.EscapeString(examples.Target)

	return &BaseHandler{
		db:       db,
		kvStore:  kvStore,
		config:   config,
		logger:   logger,
		temporal: temporal,
		store:    store,
		examples: examples,
	}
}
