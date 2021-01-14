package migrations

import (
	"github.com/baking-bad/bcdhub/internal/config"
	"github.com/baking-bad/bcdhub/internal/logger"
	"github.com/baking-bad/bcdhub/internal/metrics"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/models/transfer"
	transferParsers "github.com/baking-bad/bcdhub/internal/parsers/transfer"
)

// InitialStorageEvents -
type InitialStorageEvents struct {
	contracts map[string]string
}

// Key -
func (m *InitialStorageEvents) Key() string {
	return "execute_initial_storage"
}

// Description -
func (m *InitialStorageEvents) Description() string {
	return "execute all initial storages"
}

// Do - migrate function
func (m *InitialStorageEvents) Do(ctx *config.Context) error {
	m.contracts = make(map[string]string)
	tzips, err := ctx.TZIP.GetWithEvents()
	if err != nil {
		return err
	}

	logger.Info("Found %d tzips", len(tzips))

	h := metrics.New(ctx.Contracts, ctx.BigMapDiffs, ctx.Blocks, ctx.Protocols, ctx.Operations, ctx.Schema, ctx.TokenBalances, ctx.TZIP, ctx.Migrations, ctx.Storage, ctx.Bulk, ctx.DB)

	logger.Info("Execution events...")
	newTransfers := make([]*transfer.Transfer, 0)
	for i := range tzips {
		logger.Info("%s...", tzips[i].Address)

		rpc, err := ctx.GetRPC(tzips[i].Network)
		if err != nil {
			return err
		}
		transfers, err := h.ExecuteInitialStorageEvent(rpc, &tzips[i])
		if err != nil {
			return err
		}
		for i := range transfers {
			found, err := ctx.Transfers.Get(transfer.GetContext{
				Hash:    transfers[i].Hash,
				Network: transfers[i].Network,
				TokenID: -1,
			})
			if err != nil {
				if !ctx.Storage.IsRecordNotFound(err) {
					return err
				}
			}
			if len(found.Transfers) > 0 {
				continue
			}

			newTransfers = append(newTransfers, transfers[i])
			m.contracts[transfers[i].Contract] = transfers[i].Network
		}
	}

	updated := make([]models.Model, 0)
	if len(newTransfers) == 0 {
		return nil
	}
	for i := range newTransfers {
		updated = append(updated, newTransfers[i])
	}
	logger.Info("Found %d transfers", len(updated))
	if err := ctx.Bulk.Insert(updated); err != nil {
		return err
	}
	return transferParsers.UpdateTokenBalances(ctx.TokenBalances, newTransfers)
}

// AffectedContracts -
func (m *InitialStorageEvents) AffectedContracts() map[string]string {
	return m.contracts
}