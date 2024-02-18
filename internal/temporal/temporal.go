package temporal

import (
	"context"
	"time"

	"code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/jwt"
	"code.tjo.space/mentos1386/zdravko/pkg/retry"
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

func ConnectServerToTemporal(cfg *config.Config) (client.Client, error) {
	// For server we generate new token with admin permissions
	token, err := jwt.NewToken(cfg, []string{"temporal-system:admin", "default:admin"}, "server")
	if err != nil {
		return nil, err
	}

	provider := &AuthHeadersProvider{token}

	// Try to connect to the Temporal Server
	return retry.Retry(5, 6*time.Second, func() (client.Client, error) {
		return client.Dial(client.Options{
			HostPort:        cfg.Temporal.ServerHost,
			HeadersProvider: provider,
		})
	})
}

func ConnectWorkerToTemporal(cfg *config.Config) (client.Client, error) {
	provider := &AuthHeadersProvider{cfg.Worker.Token}

	// Try to connect to the Temporal Server
	return retry.Retry(5, 6*time.Second, func() (client.Client, error) {
		return client.Dial(client.Options{
			HostPort:        cfg.Temporal.ServerHost,
			HeadersProvider: provider,
			Namespace:       "default",
		})
	})
}
