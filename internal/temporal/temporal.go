package temporal

import (
	"context"
	"log/slog"
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/jwt"
	"code.tjo.space/mentos1386/zdravko/pkg/retry"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/client"
)

type AuthHeadersProvider struct {
	Token string
}

func (p *AuthHeadersProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + p.Token,
	}, nil
}

func ConnectServerToTemporal(logger *slog.Logger, cfg *config.ServerConfig) (client.Client, error) {
	// For server we generate new token with admin permissions
	token, err := jwt.NewTokenForServer(cfg.Jwt.PrivateKey, cfg.Jwt.PublicKey)
	if err != nil {
		return nil, err
	}

	provider := &AuthHeadersProvider{token}

	// Try to connect to the Temporal Server
	c, err := retry.Retry(10, 2*time.Second, func() (client.Client, error) {
		return client.Dial(client.Options{
			HostPort:        cfg.Temporal.ServerHost,
			HeadersProvider: provider,
		})
	})
	if err != nil {
		logger.Error("Failed to connect to Temporal Server after retries")
		return nil, errors.Wrap(err, "failed to connect to Temporal Server after retries")
	}

	return c, nil
}

func ConnectWorkerToTemporal(token string, temporalHost string) (client.Client, error) {
	provider := &AuthHeadersProvider{token}

	// Try to connect to the Temporal Server
	return retry.Retry(5, 6*time.Second, func() (client.Client, error) {
		client, err := client.Dial(client.Options{
			HostPort:        temporalHost,
			HeadersProvider: provider,
			Namespace:       "default",
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to Temporal Server: "+temporalHost)
		}
		return client, nil
	})
}
