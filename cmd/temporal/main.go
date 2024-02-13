package main

import (
	"log"

	"code.tjo.space/mentos1386/zdravko/internal"
	t "code.tjo.space/mentos1386/zdravko/pkg/temporal"
)

func backendServer(config *internal.Config) {
	serverConfig := t.NewServerConfig(config)

	server, err := t.NewServer(serverConfig)
	if err != nil {
		log.Fatalf("Unable to create server: %v", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}

	err = server.Stop()
	if err != nil {
		log.Fatalf("Unable to stop server: %v", err)
	}
}

func frontendServer(config *internal.Config) {
	uiConfig := t.NewUiConfig(config)

	uiServer, err := t.NewUiServer(uiConfig)
	if err != nil {
		log.Fatalf("Unable to create UI server: %v", err)
	}

	err = uiServer.Start()
	if err != nil {
		log.Fatalf("Unable to start UI server: %v", err)
	}

	uiServer.Stop()
}

func main() {
	config := internal.NewConfig()

	go func() {
		frontendServer(config)
	}()
	backendServer(config)
}
