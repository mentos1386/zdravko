package internal

import (
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"code.tjo.space/mentos1386/zdravko/internal/models/query"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:generate just _generate-gorm
func ConnectToDatabase(path string) (*gorm.DB, *query.Query, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	err = db.AutoMigrate(
		models.Worker{},
		models.Healthcheck{},
		models.HealthcheckHistory{},
		models.Cronjob{},
		models.CronjobHistory{},
		models.OAuth2State{},
	)
	if err != nil {
		return nil, nil, err
	}

	q := query.Use(db)

	return db, q, nil
}
