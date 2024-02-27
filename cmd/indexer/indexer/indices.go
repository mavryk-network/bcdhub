package indexer

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/ticket"
	"github.com/rs/zerolog/log"
)

func (bi *BlockchainIndexer) createIndices(ctx context.Context) error {
	log.Info().Str("network", bi.Network.String()).Msg("creating database indices...")

	// Big map action
	action := (*bigmapaction.BigMapAction)(nil)
	if err := bi.Storage.CreateIndex(ctx, "big_map_action_level_idx", "level", action); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "big_map_actions_source_ptr_idx", "source_ptr", action); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "big_map_actions_destination_ptr_idx", "destination_ptr", action); err != nil {
		return err
	}

	// Big map diff
	diff := (*bigmapdiff.BigMapDiff)(nil)
	if err := bi.Storage.CreateIndex(ctx, "big_map_diff_operation_id_idx", "operation_id", diff); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "big_map_diff_level_idx", "level", diff); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "big_map_diff_protocol_idx", "protocol_id", diff); err != nil {
		return err
	}

	// Contracts
	contractModel := (*contract.Contract)(nil)
	if err := bi.Storage.CreateIndex(ctx, "contracts_level_idx", "level", contractModel); err != nil {
		return err
	}

	// Global constants
	globalConstant := (*contract.GlobalConstant)(nil)
	if err := bi.Storage.CreateIndex(ctx, "global_constants_address_idx", "address", globalConstant); err != nil {
		return err
	}

	// Migrations
	migration := (*migration.Migration)(nil)
	if err := bi.Storage.CreateIndex(ctx, "migrations_level_idx", "level", migration); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "migrations_contract_id_idx", "contract_id", migration); err != nil {
		return err
	}

	// Operations
	operation := (*operation.Operation)(nil)
	if err := bi.Storage.CreateIndex(ctx, "operations_level_idx", "level", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_source_idx", "source_id", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_opg_idx", "hash, counter, content_index", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_entrypoint_idx", "entrypoint", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_hash_idx", "hash", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_opg_with_nonce_idx", "hash, counter, nonce", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_sort_idx", "level, counter, id", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_timestamp_idx", "timestamp", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_kind_idx", "kind", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_destination_timestamp_idx", "destination_id, timestamp", operation); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "operations_source_timestamp_idx", "source_id, timestamp", operation); err != nil {
		return err
	}

	// Scripts to global constants
	scriptConstants := (*contract.ScriptConstants)(nil)
	if err := bi.Storage.CreateIndex(ctx, "script_id_idx", "script_id", scriptConstants); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "global_constant_id_idx", "global_constant_id", scriptConstants); err != nil {
		return err
	}

	// Ticket updates
	ticketUpdate := (*ticket.TicketUpdate)(nil)
	if err := bi.Storage.CreateIndex(ctx, "ticket_updates_level_idx", "level", ticketUpdate); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "ticket_updates_operation_id_idx", "operation_id", ticketUpdate); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "ticket_updates_ticket_id_idx", "ticket_id", ticketUpdate); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "ticket_updates_account_id_idx", "account_id", ticketUpdate); err != nil {
		return err
	}

	// Tickets
	tickets := (*ticket.Ticket)(nil)
	if err := bi.Storage.CreateIndex(ctx, "ticket_level_idx", "level", tickets); err != nil {
		return err
	}
	if err := bi.Storage.CreateIndex(ctx, "ticket_ticketer_id_idx", "ticketer_id", tickets); err != nil {
		return err
	}

	log.Info().Str("network", bi.Network.String()).Msg("database indices was created")

	return nil
}
