package operations

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/bcd/ast"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Migration -
type Migration struct {
	contracts contract.Repository
}

// NewMigration -
func NewMigration(contracts contract.Repository) Migration {
	return Migration{contracts}
}

// Parse -
func (m Migration) Parse(ctx context.Context, data noderpc.Operation, operation *operation.Operation, protocol string, store parsers.Store) error {
	switch protocol {
	case
		"ProtoGenesisGenesisGenesisGenesisGenesisGenesk612im",
		"ProtoDemoNoopsDemoNoopsDemoNoopsDemoNoopsDemo6XBoYp",
		"Ps9mPmXaRzmzk35gbAYNCAw6UXdE2qoABTHbN2oEEc1qM7CwT9P":
		return m.fromBigMapDiffs(ctx, data, operation, store)
	case
		"ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK",
		"PtAtLasjh71tv2N8SDMtjajR42wTSAd9xFTvXvhDuYfRJPRLSL2":
		return m.fromLazyStorageDiff(ctx, data, operation, store)
	default:
		return errors.Errorf("unknown protocol for migration parser: %s", protocol)
	}
}

func (m Migration) fromLazyStorageDiff(ctx context.Context, data noderpc.Operation, operation *operation.Operation, store parsers.Store) error {
	var lsd []noderpc.LazyStorageDiff
	switch {
	case data.Result != nil && data.Result.LazyStorageDiff != nil:
		lsd = data.Result.LazyStorageDiff
	case data.Metadata != nil && data.Metadata.OperationResult != nil && data.Metadata.OperationResult.LazyStorageDiff != nil:
		lsd = data.Metadata.OperationResult.LazyStorageDiff
	default:
		return nil
	}

	for i := range lsd {
		if lsd[i].Kind != types.LazyStorageDiffBigMap || lsd[i].Diff == nil || lsd[i].Diff.BigMap == nil {
			continue
		}

		if lsd[i].Diff.BigMap.Action != types.BigMapActionStringUpdate {
			continue
		}

		for j := range lsd[i].Diff.BigMap.Updates {
			migration, err := m.createMigration(ctx, lsd[i].Diff.BigMap.Updates[j].Value, operation)
			if err != nil {
				return err
			}
			if migration != nil {
				operation.Destination.MigrationsCount += 1
				store.AddMigrations(migration)
				log.Info().Fields(migration.LogFields()).Msg("Migration detected")
			}
		}
	}
	return nil
}

func (m Migration) fromBigMapDiffs(ctx context.Context, data noderpc.Operation, operation *operation.Operation, store parsers.Store) error {
	var bmd []noderpc.BigMapDiff
	switch {
	case data.Result != nil && data.Result.BigMapDiffs != nil:
		bmd = data.Result.BigMapDiffs
	case data.Metadata != nil && data.Metadata.OperationResult != nil && data.Metadata.OperationResult.BigMapDiffs != nil:
		bmd = data.Metadata.OperationResult.BigMapDiffs
	default:
		return nil
	}

	for i := range bmd {
		if bmd[i].Action != types.BigMapActionStringUpdate {
			continue
		}

		migration, err := m.createMigration(ctx, bmd[i].Value, operation)
		if err != nil {
			return err
		}
		if migration != nil {
			operation.Destination.MigrationsCount += 1
			store.AddMigrations(migration)
			log.Info().Fields(migration.LogFields()).Msg("Migration detected")
		}
	}
	return nil
}

func (m Migration) createMigration(ctx context.Context, value []byte, operation *operation.Operation) (*migration.Migration, error) {
	if len(value) == 0 {
		return nil, nil
	}
	var tree ast.UntypedAST
	if err := json.Unmarshal(value, &tree); err != nil {
		return nil, err
	}

	if len(tree) == 0 {
		return nil, nil
	}

	if !tree[0].IsLambda() {
		return nil, nil
	}

	c, err := m.contracts.Get(ctx, operation.Destination.Address)
	if err != nil {
		return nil, err
	}
	return &migration.Migration{
		ContractID: c.ID,
		Contract:   c,
		Level:      operation.Level,
		ProtocolID: operation.ProtocolID,
		Timestamp:  operation.Timestamp,
		Hash:       operation.Hash,
		Kind:       types.MigrationKindLambda,
	}, nil
}
