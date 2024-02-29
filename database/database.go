package database

import (
	"embed"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sqlite/migrations/*
var sqliteMigrations embed.FS

func ConnectToDatabase(logger *slog.Logger, path string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", fmt.Sprintf("%s?_journal=WAL&_timeout=5000&_fk=true", path))
	if err != nil {
		return nil, err
	}
	logger.Info("Connected to database")

	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: sqliteMigrations,
		Root:       "sqlite/migrations",
	}
	n, err := migrate.Exec(db.DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return nil, err
	}
	logger.Info("Applied migrations", "number", n)

	return db, nil
}
