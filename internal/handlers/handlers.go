package handlers

import (
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/gorilla/sessions"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
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
	db     *gorm.DB
	query  *query.Query
	config *config.ServerConfig

	temporal client.Client

	store *sessions.CookieStore
}

func NewBaseHandler(db *gorm.DB, q *query.Query, temporal client.Client, config *config.ServerConfig) *BaseHandler {
	store := sessions.NewCookieStore([]byte(config.SessionSecret))

	return &BaseHandler{db, q, config, temporal, store}
}
