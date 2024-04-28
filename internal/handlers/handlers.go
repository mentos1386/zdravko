package handlers

import (
	"log/slog"

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/kv"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"go.temporal.io/sdk/client"
)

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
	kvStore kv.KeyValueStore
	config  *config.ServerConfig
	logger  *slog.Logger

	temporal client.Client

	store *sessions.CookieStore
}

func NewBaseHandler(db *sqlx.DB, kvStore kv.KeyValueStore, temporal client.Client, config *config.ServerConfig, logger *slog.Logger) *BaseHandler {
	store := sessions.NewCookieStore([]byte(config.SessionSecret))

	return &BaseHandler{
		db:       db,
		kvStore:  kvStore,
		config:   config,
		logger:   logger,
		temporal: temporal,
		store:    store,
	}
}
