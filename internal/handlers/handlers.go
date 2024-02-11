package handlers

import (
	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"gorm.io/gorm"
)

type BaseHandler struct {
	db     *gorm.DB
	query  *query.Query
	config *internal.Config
}

func NewBaseHandler(db *gorm.DB, q *query.Query, config *internal.Config) *BaseHandler {
	return &BaseHandler{db, q, config}
}
