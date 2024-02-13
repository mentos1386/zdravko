package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"go.temporal.io/server/common/config"
	"go.temporal.io/server/schema/sqlite"
	"go.temporal.io/server/temporal"

	t "code.tjo.space/mentos1386/zdravko/pkg/temporal"

	"go.temporal.io/server/common/authorization"
	tlog "go.temporal.io/server/common/log"

	uiserver "github.com/temporalio/ui-server/v2/server"
	uiconfig "github.com/temporalio/ui-server/v2/server/config"
	uiserveroptions "github.com/temporalio/ui-server/v2/server/server_options"
)

func backendServer() {
	cfg := t.NewConfig()

	logger := tlog.NewZapLogger(tlog.BuildZapLogger(tlog.Config{
		Stdout:     true,
		Level:      "info",
		OutputFile: "",
	}))

	sqlConfig := cfg.Persistence.DataStores[t.PersistenceStoreName].SQL

	// Apply migrations if file does not already exist
	if _, err := os.Stat(sqlConfig.DatabaseName); os.IsNotExist(err) {
		// Check if any of the parent dirs are missing
		dir := filepath.Dir(sqlConfig.DatabaseName)
		if _, err := os.Stat(dir); err != nil {
			log.Fatal(err)
		}

		if err := sqlite.SetupSchema(sqlConfig); err != nil {
			log.Fatal(err)
		}
	}

	// Pre-create namespaces
	var namespaces []*sqlite.NamespaceConfig
	for _, ns := range []string{"default"} {
		namespaces = append(namespaces, sqlite.NewNamespaceConfig(cfg.ClusterMetadata.CurrentClusterName, ns, false))
	}
	if err := sqlite.CreateNamespaces(sqlConfig, namespaces...); err != nil {
		log.Fatal(err)
	}

	authorizer, err := authorization.GetAuthorizerFromConfig(&cfg.Global.Authorization)
	if err != nil {
		log.Fatal(err)
	}

	claimMapper, err := authorization.GetClaimMapperFromConfig(&cfg.Global.Authorization, logger)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	interruptChan := make(chan interface{}, 1)
	go func() {
		if doneChan := ctx.Done(); doneChan != nil {
			s := <-doneChan
			interruptChan <- s
		} else {
			s := <-temporal.InterruptCh()
			interruptChan <- s
		}
	}()

	temporal, err := temporal.NewServer(
		temporal.WithConfig(cfg),
		temporal.ForServices(temporal.DefaultServices),
		temporal.WithLogger(logger),
		temporal.WithAuthorizer(authorizer),
		temporal.WithClaimMapper(func(cfg *config.Config) authorization.ClaimMapper {
			return claimMapper
		}),
		temporal.InterruptOn(interruptChan),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting temporal server")
	if err := temporal.Start(); err != nil {
		panic(err)
	}
	err = temporal.Stop()
	if err != nil {
		panic(err)
	}

}

func frontendServer() {
	cfg := &uiconfig.Config{
		Host:                "0.0.0.0",
		Port:                8223,
		TemporalGRPCAddress: "localhost:7233",
		EnableUI:            true,
		UIAssetPath:         "",
		Codec: uiconfig.Codec{
			Endpoint: "",
		},
		CORS: uiconfig.CORS{
			CookieInsecure: true,
		},
	}

	server := uiserver.NewServer(uiserveroptions.WithConfigProvider(cfg))

	log.Println("Starting temporal ui server")
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func main() {
	go func() {
		frontendServer()
	}()
	backendServer()
}
