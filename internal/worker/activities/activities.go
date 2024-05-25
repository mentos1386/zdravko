package activities

import (
	"log/slog"

	"github.com/mentos1386/zdravko/internal/config"
)

type Activities struct {
	config *config.WorkerConfig
	logger *slog.Logger
}

func NewActivities(config *config.WorkerConfig, logger *slog.Logger) *Activities {
	return &Activities{config: config, logger: logger}
}
