package internal

import (
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/pkg/retry"
	"go.temporal.io/sdk/client"
)

func ConnectToTemporal(cfg *config.Config) (client.Client, error) {
	// Try to connect to the Temporal Server
	return retry.Retry(5, 6*time.Second, func() (client.Client, error) {
		return client.Dial(client.Options{
			HostPort: cfg.Temporal.ServerHost,
		})
	})
}
