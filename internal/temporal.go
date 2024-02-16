package internal

import (
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"go.temporal.io/sdk/client"
)

func ConnectToTemporal(cfg *config.Config) (client.Client, error) {
	c, err := client.Dial(client.Options{HostPort: cfg.Temporal.ServerHost})
	if err != nil {
		return nil, err
	}
	return c, nil
}
