package main

import (
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	config := config.NewServerConfig()

	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "internal/models/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	db, err := gorm.Open(sqlite.Open(config.DatabasePath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(
		models.Worker{},
		models.Healthcheck{},
		models.HealthcheckHistory{},
		models.Cronjob{},
		models.CronjobHistory{},
		models.OAuth2State{},
	)

	// Execute the generator
	g.Execute()
}
