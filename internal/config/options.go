package config

import (
	"time"

	"github.com/mavryk-network/bcdhub/internal/bcd/tezerrors"
	"github.com/mavryk-network/bcdhub/internal/postgres/account"
	"github.com/mavryk-network/bcdhub/internal/postgres/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/postgres/contract"
	"github.com/mavryk-network/bcdhub/internal/postgres/domains"
	"github.com/mavryk-network/bcdhub/internal/postgres/global_constant"
	"github.com/mavryk-network/bcdhub/internal/postgres/migration"
	"github.com/mavryk-network/bcdhub/internal/postgres/operation"
	"github.com/mavryk-network/bcdhub/internal/postgres/protocol"
	smartrollup "github.com/mavryk-network/bcdhub/internal/postgres/smart_rollup"
	"github.com/mavryk-network/bcdhub/internal/postgres/stats"
	"github.com/mavryk-network/bcdhub/internal/postgres/ticket"
	"github.com/mavryk-network/bcdhub/internal/services/mempool"

	"github.com/mavryk-network/bcdhub/internal/postgres/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/postgres/block"
	pgCore "github.com/mavryk-network/bcdhub/internal/postgres/core"

	"github.com/mavryk-network/bcdhub/internal/noderpc"
)

// ContextOption -
type ContextOption func(ctx *Context)

// WithRPC -
func WithRPC(rpcConfig map[string]RPCConfig) ContextOption {
	return func(ctx *Context) {
		if rpcProvider, ok := rpcConfig[ctx.Network.String()]; ok {
			if rpcProvider.URI == "" {
				return
			}
			opts := []noderpc.NodeOption{
				noderpc.WithTimeout(time.Second * time.Duration(rpcProvider.Timeout)),
				noderpc.WithRateLimit(rpcProvider.RequestsPerSecond),
			}
			if rpcProvider.Log {
				opts = append(opts, noderpc.WithLog())
			}

			ctx.RPC = noderpc.NewNodeRPC(rpcProvider.URI, opts...)
		}
	}
}

// WithWaitRPC -
func WithWaitRPC(rpcConfig map[string]RPCConfig) ContextOption {
	return func(ctx *Context) {
		if rpcProvider, ok := rpcConfig[ctx.Network.String()]; ok {
			if rpcProvider.URI == "" {
				return
			}
			opts := []noderpc.NodeOption{
				noderpc.WithTimeout(time.Second * time.Duration(rpcProvider.Timeout)),
				noderpc.WithRateLimit(rpcProvider.RequestsPerSecond),
			}
			if rpcProvider.Log {
				opts = append(opts, noderpc.WithLog())
			}

			ctx.RPC = noderpc.NewWaitNodeRPC(rpcProvider.URI, opts...)
		}
	}
}

// WithStorage -
func WithStorage(cfg StorageConfig, appName string, maxPageSize int64) ContextOption {
	return func(ctx *Context) {
		if len(cfg.Postgres.Host) == 0 {
			panic("Please set connection strings to storage in config")
		}

		opts := []pgCore.PostgresOption{
			pgCore.WithPageSize(maxPageSize),
		}

		if cfg.LogQueries {
			opts = append(opts, pgCore.WithQueryLogging())
		}

		conn := pgCore.WaitNew(
			cfg.Postgres, ctx.Network.String(),
			appName, cfg.Timeout, opts...,
		)

		contractStorage := contract.NewStorage(conn)
		ctx.StorageDB = conn
		ctx.Storage = conn
		ctx.Accounts = account.NewStorage(conn)
		ctx.BigMapActions = bigmapaction.NewStorage(conn)
		ctx.Blocks = block.NewStorage(conn)
		ctx.BigMapDiffs = bigmapdiff.NewStorage(conn)
		ctx.Contracts = contractStorage
		ctx.Migrations = migration.NewStorage(conn)
		ctx.Operations = operation.NewStorage(conn)
		ctx.Protocols = protocol.NewStorage(conn)
		ctx.GlobalConstants = global_constant.NewStorage(conn)
		ctx.Domains = domains.NewStorage(conn)
		ctx.Tickets = ticket.NewStorage(conn)
		ctx.Scripts = contractStorage
		ctx.SmartRollups = smartrollup.NewStorage(conn)
		ctx.Stats = stats.NewStorage(conn)
	}
}

// WithConfigCopy -
func WithConfigCopy(cfg Config) ContextOption {
	return func(ctx *Context) {
		ctx.Config = cfg
	}
}

// WithMempool -
func WithMempool(cfg map[string]ServiceConfig) ContextOption {
	return func(ctx *Context) {
		if svcCfg, ok := cfg[ctx.Network.String()]; ok {
			if svcCfg.MempoolURI == "" {
				return
			}
			ctx.Mempool = mempool.NewMempool(svcCfg.MempoolURI)
		}
	}
}

// WithLoadErrorDescriptions -
func WithLoadErrorDescriptions() ContextOption {
	return func(ctx *Context) {
		if err := tezerrors.LoadErrorDescriptions(); err != nil {
			panic(err)
		}
	}
}
