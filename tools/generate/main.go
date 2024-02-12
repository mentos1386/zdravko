package main

import (
	"code.tjo.space/mentos1386/zdravko/internal"
	"code.tjo.space/mentos1386/zdravko/internal/models"
	"gorm.io/gen"
)

func main() {
	config := internal.NewConfig()

	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "internal/models/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	db, _, _ := internal.ConnectToDatabase(config.SQLITE_DB_PATH)

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(models.Healthcheck{}, models.OAuth2State{})

	// Execute the generator
	g.Execute()
}
