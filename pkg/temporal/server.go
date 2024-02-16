package temporal

import (
	"context"
	"os"
	"path/filepath"

	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/schema/sqlite"
	t "go.temporal.io/server/temporal"
)

func NewServer(cfg *config.Config) (t.Server, error) {
	logger := log.NewZapLogger(log.BuildZapLogger(log.Config{
		Stdout:     true,
		Level:      "info",
		OutputFile: "",
	}))

	sqlConfig := cfg.Persistence.DataStores[PersistenceStoreName].SQL

	// Apply migrations if file does not already exist
	if _, err := os.Stat(sqlConfig.DatabaseName); os.IsNotExist(err) {
		// Check if any of the parent dirs are missing
		dir := filepath.Dir(sqlConfig.DatabaseName)
		if _, err := os.Stat(dir); err != nil {
			return nil, err
		}

		if err := sqlite.SetupSchema(sqlConfig); err != nil {
			return nil, err
		}
	}

	// Pre-create namespaces
	var namespaces []*sqlite.NamespaceConfig
	for _, ns := range []string{"default"} {
		namespaces = append(namespaces, sqlite.NewNamespaceConfig(cfg.ClusterMetadata.CurrentClusterName, ns, false))
	}
	if err := sqlite.CreateNamespaces(sqlConfig, namespaces...); err != nil {
		return nil, err
	}

	authorizer, err := authorization.GetAuthorizerFromConfig(&cfg.Global.Authorization)
	if err != nil {
		return nil, err
	}

	claimMapper, err := authorization.GetClaimMapperFromConfig(&cfg.Global.Authorization, logger)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	interruptChan := make(chan interface{}, 1)
	go func() {
		if doneChan := ctx.Done(); doneChan != nil {
			s := <-doneChan
			interruptChan <- s
		} else {
			s := <-t.InterruptCh()
			interruptChan <- s
		}
	}()

	return t.NewServer(
		t.WithConfig(cfg),
		t.ForServices(t.DefaultServices),
		t.WithLogger(logger),
		t.WithAuthorizer(authorizer),
		t.WithClaimMapper(func(cfg *config.Config) authorization.ClaimMapper {
			return claimMapper
		}),
		t.InterruptOn(interruptChan),
	)
}