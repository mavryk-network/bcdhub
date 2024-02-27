package store

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/ticket"
	"github.com/mavryk-network/bcdhub/internal/postgres/core"
	"github.com/pkg/errors"
)

// Save -
func (store *Store) Save(ctx context.Context) error {
	stats, err := store.stats.Get(ctx)
	if err != nil {
		return err
	}
	store.Stats.ID = stats.ID

	tx, err := core.NewTransaction(ctx, store.db)
	if err != nil {
		return err
	}

	if err := tx.Block(ctx, store.Block); err != nil {
		return errors.Wrap(err, "saving block")
	}

	if err := store.saveAccounts(ctx, tx); err != nil {
		return errors.Wrap(err, "saving accounts")
	}

	if err := store.saveTickets(ctx, tx); err != nil {
		return errors.Wrap(err, "saving tickets")
	}

	if err := store.saveTicketBalances(ctx, tx); err != nil {
		return errors.Wrap(err, "saving ticket balances")
	}

	if err := store.saveOperations(ctx, tx); err != nil {
		return errors.Wrap(err, "saving operations")
	}

	if err := store.saveContracts(ctx, tx); err != nil {
		return errors.Wrap(err, "saving contracts")
	}

	if err := store.saveMigrations(ctx, tx); err != nil {
		return errors.Wrap(err, "saving migrations")
	}

	if err := tx.BigMapStates(ctx, store.bigMapStates()...); err != nil {
		return errors.Wrap(err, "saving bigmap states")
	}

	if err := tx.GlobalConstants(ctx, store.GlobalConstants...); err != nil {
		return errors.Wrap(err, "saving bigmap states")
	}

	if err := store.saveSmartRollups(ctx, tx); err != nil {
		return errors.Wrap(err, "saving smart rollups")
	}

	if err := tx.UpdateStats(ctx, store.Stats); err != nil {
		return errors.Wrap(err, "saving stats")
	}

	return tx.Commit()
}

func (store *Store) saveAccounts(ctx context.Context, tx models.Transaction) error {
	if len(store.Accounts) == 0 {
		return nil
	}

	arr := make([]*account.Account, 0, len(store.Accounts))
	for _, acc := range store.Accounts {
		if acc.IsEmpty() {
			continue
		}
		arr = append(arr, acc)
	}

	if err := tx.Accounts(ctx, arr...); err != nil {
		return err
	}

	for i := range arr {
		store.accIds[arr[i].Address] = arr[i].ID
	}

	return nil
}

func (store *Store) saveTickets(ctx context.Context, tx models.Transaction) error {
	if len(store.Tickets) == 0 {
		return nil
	}

	arr := make([]*ticket.Ticket, 0, len(store.Tickets))
	for _, t := range store.Tickets {
		if id, ok := store.getAccountId(t.Ticketer); ok {
			t.TicketerID = id
		} else {
			return errors.Errorf("unknown ticketer account: %s", t.Ticketer.Address)
		}
		arr = append(arr, t)
	}

	if err := tx.Tickets(ctx, arr...); err != nil {
		return err
	}

	for i := range arr {
		store.ticketIds[arr[i].GetHash()] = arr[i].ID
	}

	return nil
}

func (store *Store) saveTicketBalances(ctx context.Context, tx models.Transaction) error {
	if len(store.TicketBalances) == 0 {
		return nil
	}

	balances := make([]*ticket.Balance, 0, len(store.TicketBalances))
	for _, balance := range store.TicketBalances {
		if id, ok := store.getAccountId(balance.Account); ok {
			balance.AccountId = id
		} else {
			return errors.Errorf("unknown ticket balance account: %s", balance.Account.Address)
		}
		if id, ok := store.getAccountId(balance.Ticket.Ticketer); ok {
			balance.Ticket.TicketerID = id
		} else {
			return errors.Errorf("unknown ticket balance ticketer: %s", balance.Ticket.Ticketer.Address)
		}

		if id, ok := store.ticketIds[balance.Ticket.GetHash()]; ok {
			balance.TicketId = id
		} else {
			return errors.Errorf("unknown ticket of balance: %s", balance.Ticket.Ticketer.Address)
		}
		balances = append(balances, balance)
	}

	return tx.TicketBalances(ctx, balances...)
}

func (store *Store) saveMigrations(ctx context.Context, tx models.Transaction) error {
	if len(store.Migrations) == 0 {
		return nil
	}

	for i := range store.Migrations {
		if store.Migrations[i].ContractID == 0 {
			store.Migrations[i].ContractID = store.Migrations[i].Contract.ID
		}
	}

	return tx.Migrations(ctx, store.Migrations...)
}

func (store *Store) saveSmartRollups(ctx context.Context, tx models.Transaction) error {
	if len(store.SmartRollups) == 0 {
		return nil
	}

	for i, rollup := range store.SmartRollups {
		if id, ok := store.getAccountId(rollup.Address); ok {
			store.SmartRollups[i].AddressId = id
		} else {
			return errors.Errorf("unknown smart rollup account: %s", rollup.Address.Address)
		}
	}

	return tx.SmartRollups(ctx, store.SmartRollups...)
}

func (store *Store) saveOperations(ctx context.Context, tx models.Transaction) error {
	if len(store.Operations) == 0 {
		return nil
	}

	for i := range store.Operations {
		if err := store.setOperationAccountsId(store.Operations[i]); err != nil {
			return err
		}
	}

	if err := tx.Operations(ctx, store.Operations...); err != nil {
		return errors.Wrap(err, "saving operations")
	}

	var (
		bigMapDiffs   = make([]*bigmapdiff.BigMapDiff, 0)
		bigMapActions = make([]*bigmapaction.BigMapAction, 0)
		ticketUpdates = make([]*ticket.TicketUpdate, 0)
	)

	for _, operation := range store.Operations {
		for j := range operation.BigMapDiffs {
			operation.BigMapDiffs[j].OperationID = operation.ID
		}
		bigMapDiffs = append(bigMapDiffs, operation.BigMapDiffs...)

		for j := range operation.BigMapActions {
			operation.BigMapActions[j].OperationID = operation.ID
		}
		bigMapActions = append(bigMapActions, operation.BigMapActions...)

		for j, update := range operation.TicketUpdates {
			if id, ok := store.getAccountId(update.Account); ok {
				operation.TicketUpdates[j].AccountId = id
			} else {
				return errors.Errorf("unknown ticket update account: %s", update.Account.Address)
			}

			if id, ok := store.getAccountId(update.Ticket.Ticketer); ok {
				operation.TicketUpdates[j].Ticket.TicketerID = id
			} else {
				return errors.Errorf("unknown ticket update ticketer account: %s", update.Ticket.Ticketer.Address)
			}

			operation.TicketUpdates[j].OperationId = operation.ID

			hash := operation.TicketUpdates[j].Ticket.GetHash()
			if id, ok := store.ticketIds[hash]; ok {
				operation.TicketUpdates[j].TicketId = id
			} else {
				return errors.Errorf("unknown ticket: ticketer_id=%d content_type=%s content=%s",
					operation.TicketUpdates[j].Ticket.TicketerID,
					operation.TicketUpdates[j].Ticket.ContentType,
					operation.TicketUpdates[j].Ticket.Content,
				)
			}
		}

		ticketUpdates = append(ticketUpdates, operation.TicketUpdates...)
	}

	if err := tx.BigMapDiffs(ctx, bigMapDiffs...); err != nil {
		return errors.Wrap(err, "saving bigmap diffs")
	}
	if err := tx.BigMapActions(ctx, bigMapActions...); err != nil {
		return errors.Wrap(err, "saving bigmap actions")
	}
	if err := tx.TickerUpdates(ctx, ticketUpdates...); err != nil {
		return errors.Wrap(err, "saving ticket updates")
	}
	return nil
}

func (store *Store) setOperationAccountsId(operation *operation.Operation) error {
	if id, ok := store.getAccountId(operation.Source); ok {
		operation.SourceID = id
	} else {
		return errors.Errorf("unknown source account: %s", operation.Source.Address)
	}

	if !(operation.IsOrigination() && !operation.IsApplied()) {
		if id, ok := store.getAccountId(operation.Destination); ok {
			operation.DestinationID = id
		} else {
			return errors.Errorf("unknown destination account: %s", operation.Destination.Address)
		}
	}

	if id, ok := store.getAccountId(operation.Initiator); ok {
		operation.InitiatorID = id
	} else {
		return errors.Errorf("unknown initiator account: %s", operation.Initiator.Address)
	}

	if id, ok := store.getAccountId(operation.Delegate); ok {
		operation.DelegateID = id
	} else {
		return errors.Errorf("unknown delegate account: %s", operation.Delegate.Address)
	}

	return nil
}

func (store *Store) getAccountId(acc account.Account) (int64, bool) {
	if acc.IsEmpty() {
		return 0, true
	}
	id, ok := store.accIds[acc.Address]
	return id, ok
}

func (store *Store) saveContracts(ctx context.Context, tx models.Transaction) error {
	if len(store.Contracts) == 0 {
		return nil
	}

	for i := range store.Contracts {
		if store.Contracts[i].Alpha.Code != nil {
			if err := tx.Scripts(ctx, &store.Contracts[i].Alpha); err != nil {
				return err
			}
			store.Contracts[i].AlphaID = store.Contracts[i].Alpha.ID
		}
		if store.Contracts[i].Babylon.Code != nil {
			if store.Contracts[i].Alpha.Hash != store.Contracts[i].Babylon.Hash {
				if err := tx.Scripts(ctx, &store.Contracts[i].Babylon); err != nil {
					return err
				}
				store.Contracts[i].BabylonID = store.Contracts[i].Babylon.ID

				if len(store.Contracts[i].Babylon.Constants) > 0 {
					for j := range store.Contracts[i].Babylon.Constants {
						relation := contract.ScriptConstants{
							ScriptId:         store.Contracts[i].BabylonID,
							GlobalConstantId: store.Contracts[i].Babylon.Constants[j].ID,
						}
						if err := tx.ScriptConstant(ctx, &relation); err != nil {
							return err
						}
					}
				}

			} else {
				store.Contracts[i].BabylonID = store.Contracts[i].Alpha.ID
			}
		}
		if store.Contracts[i].Jakarta.Code != nil {
			if store.Contracts[i].Babylon.Hash != store.Contracts[i].Jakarta.Hash {
				if err := tx.Scripts(ctx, &store.Contracts[i].Jakarta); err != nil {
					return err
				}
				store.Contracts[i].JakartaID = store.Contracts[i].Jakarta.ID

				if len(store.Contracts[i].Jakarta.Constants) > 0 {
					for j := range store.Contracts[i].Jakarta.Constants {
						relation := contract.ScriptConstants{
							ScriptId:         store.Contracts[i].JakartaID,
							GlobalConstantId: store.Contracts[i].Jakarta.Constants[j].ID,
						}
						if err := tx.ScriptConstant(ctx, &relation); err != nil {
							return err
						}
					}
				}

			} else {
				store.Contracts[i].JakartaID = store.Contracts[i].Babylon.ID
			}
		}

		if id, ok := store.getAccountId(store.Contracts[i].Account); ok {
			store.Contracts[i].AccountID = id
		} else {
			return errors.Errorf("unknown contract account: %s", store.Contracts[i].Account.Address)
		}

		if id, ok := store.getAccountId(store.Contracts[i].Manager); ok {
			store.Contracts[i].ManagerID = id
		} else {
			return errors.Errorf("unknown manager account: %s", store.Contracts[i].Manager.Address)
		}

		if id, ok := store.getAccountId(store.Contracts[i].Delegate); ok {
			store.Contracts[i].DelegateID = id
		} else {
			return errors.Errorf("unknown delegate account: %s", store.Contracts[i].Delegate.Address)
		}
	}

	if err := tx.Contracts(ctx, store.Contracts...); err != nil {
		return err
	}

	return nil
}
