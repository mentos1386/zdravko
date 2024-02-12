package handlers

import (
	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"code.tjo.space/mentos1386/zdravko/web/templates/components"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

var Pages = []*components.Page{
	{Path: "/", Title: "Status"},
	{Path: "/incidents", Title: "Incidents"},
	{Path: "/settings", Title: "Settings"},
}

func GetPageByTitle(title string) *components.Page {
	for _, p := range Pages {
		if p.Title == title {
			return p
		}
	}
	return nil
}

type BaseHandler struct {
	db     *gorm.DB
	query  *query.Query
	config *internal.Config

	store *sessions.CookieStore
}

func NewBaseHandler(db *gorm.DB, q *query.Query, config *internal.Config) *BaseHandler {
	store := sessions.NewCookieStore([]byte(config.SESSION_SECRET))

	return &BaseHandler{db, q, config, store}
}
