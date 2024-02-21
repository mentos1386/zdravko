package activities

import "code.tjo.space/mentos1386/zdravko/internal/config"

type Activities struct {
	config *config.WorkerConfig
}

func NewActivities(config *config.WorkerConfig) *Activities {
	return &Activities{config: config}
}
