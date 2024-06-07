package temporal

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/primitives"
	"go.temporal.io/server/schema/sqlite"
	t "go.temporal.io/server/temporal"
)

func NewServer(l *slog.Logger, cfg *config.Config, tokenKeyProvider authorization.TokenKeyProvider) (t.Server, error) {
	logger := slogLogger{
		log:   l,
		level: slog.LevelInfo,
	}

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

	authorizer := authorization.NewDefaultAuthorizer()
	claimMapper := authorization.NewDefaultJWTClaimMapper(tokenKeyProvider, &cfg.Global.Authorization, logger)

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
		t.ForServices([]string{
			string(primitives.FrontendService),
			string(primitives.HistoryService),
			string(primitives.MatchingService),
			string(primitives.WorkerService),
			string(primitives.InternalFrontendService),
		}),
		t.WithLogger(logger),
		t.InterruptOn(interruptChan),
		t.WithAuthorizer(authorizer),
		t.WithClaimMapper(func(cfg *config.Config) authorization.ClaimMapper {
			return claimMapper
		}),
	)
}
