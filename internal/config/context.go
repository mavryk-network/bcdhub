package config

import (
	"github.com/mavryk-network/bcdhub/internal/cache"
	"github.com/mavryk-network/bcdhub/internal/models"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/models/block"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/domains"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/protocol"
	smartrollup "github.com/mavryk-network/bcdhub/internal/models/smart_rollup"
	"github.com/mavryk-network/bcdhub/internal/models/stats"
	"github.com/mavryk-network/bcdhub/internal/models/ticket"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/postgres/core"
	"github.com/mavryk-network/bcdhub/internal/services/mempool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Context -
type Context struct {
	Network types.Network
	RPC     noderpc.INode
	Mempool *mempool.Mempool

	StorageDB *core.Postgres

	Config     Config
	TzipSchema string

	Storage         models.GeneralRepository
	Accounts        account.Repository
	BigMapActions   bigmapaction.Repository
	BigMapDiffs     bigmapdiff.Repository
	Blocks          block.Repository
	Contracts       contract.Repository
	GlobalConstants contract.ConstantRepository
	Migrations      migration.Repository
	Operations      operation.Repository
	Protocols       protocol.Repository
	Tickets         ticket.Repository
	Domains         domains.Repository
	Scripts         contract.ScriptRepository
	SmartRollups    smartrollup.Repository
	Stats           stats.Repository

	Cache *cache.Cache
}

// NewContext -
func NewContext(network types.Network, opts ...ContextOption) *Context {
	ctx := &Context{
		Network: network,
	}

	for _, opt := range opts {
		opt(ctx)
	}

	ctx.Cache = cache.NewCache(
		ctx.RPC, ctx.Accounts, ctx.Contracts, ctx.Protocols,
	)
	return ctx
}

// Close -
func (ctx *Context) Close() error {
	if ctx.StorageDB != nil {
		if err := ctx.StorageDB.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Contexts -
type Contexts map[types.Network]*Context

// NewContext -
func NewContexts(cfg Config, networks []string, opts ...ContextOption) Contexts {
	if len(networks) == 0 {
		panic("empty networks list in config file")
	}

	ctxs := make(Contexts)

	for i := range networks {
		networkType := types.NewNetwork(networks[i])
		if networkType == types.Empty {
			log.Warn().Str("network", networks[i]).Msg("unknown network")
			continue
		}
		ctxs[networkType] = NewContext(networkType, opts...)
	}

	return ctxs
}

// Get -
func (ctxs Contexts) Get(network types.Network) (*Context, error) {
	if ctx, ok := ctxs[network]; ok {
		return ctx, nil
	}
	return nil, errors.Errorf("unknown network: %s", network.String())
}

// Any -
func (ctxs Contexts) Any() *Context {
	for _, ctx := range ctxs {
		return ctx
	}
	panic("empty contexts map")
}

// Close -
func (ctxs Contexts) Close() error {
	for _, c := range ctxs {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}
