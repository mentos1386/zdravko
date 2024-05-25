package activities

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/mentos1386/zdravko/database"
	"github.com/mentos1386/zdravko/internal/config"
)

type Activities struct {
	config  *config.ServerConfig
	db      *sqlx.DB
	kvStore database.KeyValueStore
	logger  *slog.Logger
}

func NewActivities(config *config.ServerConfig, logger *slog.Logger, db *sqlx.DB, kvStore database.KeyValueStore) *Activities {
	return &Activities{config: config, logger: logger, db: db, kvStore: kvStore}
}
