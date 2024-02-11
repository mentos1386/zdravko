package internal

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// go:generate just _generate-gorm
func ConnectToDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
