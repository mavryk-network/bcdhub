package rollback

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/models/block"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	smartrollup "github.com/mavryk-network/bcdhub/internal/models/smart_rollup"
	"github.com/mavryk-network/bcdhub/internal/models/stats"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Manager -
type Manager struct {
	storage   models.GeneralRepository
	blockRepo block.Repository
	rollback  models.Rollback
	statsRepo stats.Repository
}

// NewManager -
func NewManager(
	storage models.GeneralRepository,
	blockRepo block.Repository,
	rollback models.Rollback,
	statsRepo stats.Repository,
) Manager {
	return Manager{
		storage:   storage,
		blockRepo: blockRepo,
		rollback:  rollback,
		statsRepo: statsRepo,
	}
}

// Rollback - rollback indexer state to level
func (rm Manager) Rollback(ctx context.Context, network types.Network, fromState block.Block, toLevel int64) error {
	if toLevel >= fromState.Level {
		return errors.Errorf("To level must be less than from level: %d >= %d", toLevel, fromState.Level)
	}

	for level := fromState.Level; level > toLevel; level-- {
		log.Info().Str("network", network.String()).Msgf("start rollback to %d", level)

		if _, err := rm.blockRepo.Get(ctx, level); err != nil {
			if rm.storage.IsRecordNotFound(err) {
				continue
			}
			return err
		}

		if err := rm.rollbackBlock(ctx, level); err != nil {
			log.Err(err).Str("network", network.String()).Msg("rollback error")
			return rm.rollback.Rollback()
		}

		log.Info().Str("network", network.String()).Msgf("rolled back to %d", level)
	}

	return rm.rollback.Commit()
}

func (rm Manager) rollbackBlock(ctx context.Context, level int64) error {
	rollbackCtx, err := newRollbackContext(ctx, rm.statsRepo)
	if err != nil {
		return err
	}

	if err := rm.rollbackOperations(ctx, level, &rollbackCtx); err != nil {
		return err
	}
	if err := rm.rollbackBigMapState(ctx, level); err != nil {
		return err
	}
	if err := rm.rollbackScripts(ctx, level); err != nil {
		return err
	}
	if err := rm.rollbackMigrations(ctx, level, &rollbackCtx); err != nil {
		return err
	}
	if err := rm.rollbackTickets(ctx, level); err != nil {
		return err
	}
	if err := rm.rollbackAll(ctx, level, &rollbackCtx); err != nil {
		return err
	}
	if err := rm.rollback.Protocols(ctx, level); err != nil {
		return err
	}
	if err := rollbackCtx.update(ctx, rm.rollback); err != nil {
		return err
	}

	return nil
}

func (rm Manager) rollbackMigrations(ctx context.Context, level int64, rCtx *rollbackContext) error {
	migrations, err := rm.rollback.GetMigrations(ctx, level)
	if err != nil {
		return nil
	}
	if len(migrations) == 0 {
		return nil
	}

	for i := range migrations {
		rCtx.applyMigration(migrations[i].Contract.AccountID)
	}

	if _, err := rm.rollback.DeleteAll(ctx, (*migration.Migration)(nil), level); err != nil {
		return err
	}
	log.Info().Msg("rollback migrations")
	return nil
}

func (rm Manager) rollbackAll(ctx context.Context, level int64, rCtx *rollbackContext) error {
	for _, model := range []models.Model{
		(*block.Block)(nil),
		(*bigmapdiff.BigMapDiff)(nil),
		(*bigmapaction.BigMapAction)(nil),
		(*smartrollup.SmartRollup)(nil),
		(*account.Account)(nil),
	} {
		if _, err := rm.rollback.DeleteAll(ctx, model, level); err != nil {
			return err
		}
		log.Info().Msgf("rollback: %T", model)
	}

	contractsCount, err := rm.rollback.DeleteAll(ctx, (*contract.Contract)(nil), level)
	if err != nil {
		return err
	}
	rCtx.generalStats.ContractsCount -= contractsCount
	log.Info().Msgf("rollback contracts")

	return nil
}
