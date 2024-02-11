package internal

import (
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

//go:generate just _generate-gorm
func ConnectToDatabase(path string) (*gorm.DB, *query.Query, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	db.AutoMigrate(&models.Healthcheck{})

	q := query.Use(db)

	return db, q, nil
}
