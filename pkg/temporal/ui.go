package temporal

import (
	"code.tjo.space/mentos1386/zdravko/internal"
	"github.com/temporalio/ui-server/v2/server"
	"github.com/temporalio/ui-server/v2/server/config"
	"github.com/temporalio/ui-server/v2/server/server_options"
)

func NewUiConfig(cfg *internal.Config) *config.Config {
	return &config.Config{
		Host:                cfg.Temporal.ListenAddress,
		Port:                8223,
		TemporalGRPCAddress: cfg.Temporal.ServerHost,
		EnableUI:            true,
		PublicPath:          "/temporal",
		Codec: config.Codec{
			Endpoint: "",
		},
		CORS: config.CORS{
			CookieInsecure: true,
		},
	}
}

func NewUiServer(cfg *config.Config) (*server.Server, error) {
	s := server.NewServer(server_options.WithConfigProvider(cfg))
	return s, nil
}
