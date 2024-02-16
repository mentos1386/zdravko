package temporal

import (
	"code.tjo.space/mentos1386/zdravko/internal/config"
	"github.com/temporalio/ui-server/v2/server"
	t "go.temporal.io/server/temporal"
)

type Temporal struct {
	server   t.Server
	uiServer *server.Server
}

func NewTemporal(cfg *config.Config) (*Temporal, error) {
	serverConfig := NewServerConfig(cfg)
	server, err := NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	uiConfig := NewUiConfig(cfg)
	uiServer, err := NewUiServer(uiConfig)
	if err != nil {
		return nil, err
	}

	return &Temporal{
		server:   server,
		uiServer: uiServer,
	}, nil
}

func (t *Temporal) Name() string {
	return "Temporal UI and Server"
}

func (t *Temporal) Start() error {
	go func() {
		err := t.uiServer.Start()
		if err != nil {
			panic(err)
		}
	}()
	return t.server.Start()
}

func (t *Temporal) Stop() error {
	err := t.server.Stop()
	if err != nil {
		return err
	}
	t.uiServer.Stop()

	return nil
}
