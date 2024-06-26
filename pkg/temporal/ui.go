package temporal

import (
	internal "github.com/mentos1386/zdravko/internal/config"
	"github.com/temporalio/ui-server/v2/server"
	"github.com/temporalio/ui-server/v2/server/config"
	"github.com/temporalio/ui-server/v2/server/server_options"
)

func NewUiConfig(cfg *internal.TemporalConfig) *config.Config {
	return &config.Config{
		Host:                cfg.ListenAddress,
		Port:                8223,
		TemporalGRPCAddress: "0.0.0.0:7233",
		EnableUI:            true,
		PublicPath:          "/settings/temporal",
	}
}

func NewUiServer(cfg *config.Config) (*server.Server, error) {
	s := server.NewServer(server_options.WithConfigProvider(cfg))
	return s, nil
}
