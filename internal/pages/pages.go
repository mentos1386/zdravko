package pages

import "gorm.io/gorm"

type PageHandler struct {
	db *gorm.DB
}

func NewPageHandler(db *gorm.DB) *PageHandler {
	return &PageHandler{db}
}
