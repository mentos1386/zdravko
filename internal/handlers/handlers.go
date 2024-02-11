package handlers

import (
	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

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
