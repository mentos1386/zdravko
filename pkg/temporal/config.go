package temporal

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"time"

	internal "code.tjo.space/mentos1386/zdravko/internal/config"
	"code.tjo.space/mentos1386/zdravko/internal/jwt"
	"go.temporal.io/server/common/cluster"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/persistence/sql/sqlplugin/sqlite"
)

const PersistenceStoreName = "sqlite-default"

const BroadcastAddress = "127.0.0.1"

const FrontendHttpPort = 8233
const FrontendPort = 7233
const HistoryPort = 7234
const MatchingPort = 7235
const WorkerPort = 7236
const InternalFrontendPort = 7237

type TokenKeyProvider struct {
	config *internal.TemporalConfig
}

func (p *TokenKeyProvider) SupportedMethods() []string {
	return []string{"RS256", "RS384", "RS512"}
}

func (p *TokenKeyProvider) HmacKey(alg string, kid string) ([]byte, error) {
	return nil, fmt.Errorf("HMAC key is not supported")
}

func (p *TokenKeyProvider) EcdsaKey(alg string, kid string) (*ecdsa.PublicKey, error) {
	return nil, fmt.Errorf("ECDSA key is not supported")
}

func (p *TokenKeyProvider) RsaKey(alg string, kid string) (*rsa.PublicKey, error) {
	return jwt.JwtPublicKey(p.config.Jwt.PublicKey)
}

func (p *TokenKeyProvider) Close() {
}

func NewServerConfig(cfg *internal.TemporalConfig) *config.Config {
	return &config.Config{
		Persistence: config.Persistence{
			DataStores: map[string]config.DataStore{
				PersistenceStoreName: {SQL: &config.SQL{
					PluginName: sqlite.PluginName,
					ConnectAttributes: map[string]string{
						"mode": "rwc",
					},
					DatabaseName: cfg.DatabasePath,
				},
				},
			},
			DefaultStore:     PersistenceStoreName,
			VisibilityStore:  PersistenceStoreName,
			NumHistoryShards: 1,
		},
		Global: config.Global{
			Membership: config.Membership{
				MaxJoinDuration:  30 * time.Second,
				BroadcastAddress: BroadcastAddress,
			},
			Authorization: config.Authorization{},
		},
		Services: map[string]config.Service{
			"frontend": {
				RPC: config.RPC{
					HTTPPort:        FrontendHttpPort,
					GRPCPort:        FrontendPort,
					MembershipPort:  FrontendPort + 100,
					BindOnLocalHost: false,
					BindOnIP:        cfg.ListenAddress,
				},
			},
			"history": {
				RPC: config.RPC{
					GRPCPort:        HistoryPort,
					MembershipPort:  HistoryPort + 100,
					BindOnLocalHost: true,
					BindOnIP:        "",
				},
			},
			"matching": {
				RPC: config.RPC{
					GRPCPort:        MatchingPort,
					MembershipPort:  MatchingPort + 100,
					BindOnLocalHost: true,
					BindOnIP:        "",
				},
			},
			"worker": {
				RPC: config.RPC{
					GRPCPort:        WorkerPort,
					MembershipPort:  WorkerPort + 100,
					BindOnLocalHost: true,
					BindOnIP:        "",
				},
			},
			"internal-frontend": {
				RPC: config.RPC{
					GRPCPort:        InternalFrontendPort,
					MembershipPort:  InternalFrontendPort + 100,
					BindOnLocalHost: true,
					BindOnIP:        "",
				},
			},
		},
		ClusterMetadata: &cluster.Config{
			EnableGlobalNamespace:    false,
			FailoverVersionIncrement: 10,
			MasterClusterName:        "active",
			CurrentClusterName:       "active",
			ClusterInformation: map[string]cluster.ClusterInformation{
				"active": {
					Enabled:                true,
					InitialFailoverVersion: 1,
					RPCAddress:             fmt.Sprintf("%s:%d", BroadcastAddress, InternalFrontendPort),
					ClusterID:              "todo",
				},
			},
		},
		DCRedirectionPolicy: config.DCRedirectionPolicy{
			Policy: "noop",
		},
		Archival: config.Archival{
			History: config.HistoryArchival{
				State:      "disabled",
				EnableRead: false,
				Provider:   nil,
			},
			Visibility: config.VisibilityArchival{
				State:      "disabled",
				EnableRead: false,
				Provider:   nil,
			},
		},
		NamespaceDefaults: config.NamespaceDefaults{
			Archival: config.ArchivalNamespaceDefaults{
				History: config.HistoryArchivalNamespaceDefaults{
					State: "disabled",
				},
				Visibility: config.VisibilityArchivalNamespaceDefaults{
					State: "disabled",
				},
			},
		},
	}
}
