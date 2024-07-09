package models

import (
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/models/block"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/protocol"
	smartrollup "github.com/mavryk-network/bcdhub/internal/models/smart_rollup"
	"github.com/mavryk-network/bcdhub/internal/models/stats"
	"github.com/mavryk-network/bcdhub/internal/models/ticket"
)

// Document names
const (
	DocAccounts        = "accounts"
	DocBigMapActions   = "big_map_actions"
	DocBigMapDiff      = "big_map_diffs"
	DocBigMapState     = "big_map_states"
	DocBlocks          = "blocks"
	DocContracts       = "contracts"
	DocGlobalConstants = "global_constants"
	DocMigrations      = "migrations"
	DocOperations      = "operations"
	DocProtocol        = "protocols"
	DocScripts         = "scripts"
	DocTicketUpdates   = "ticket_updates"
	DocTickets         = "tickets"
	DocTicketBalances  = "ticket_balances"
	DocSmartRollups    = "smart_rollup"
	DocStats           = "stats"
)

// AllDocuments - returns all document names
func AllDocuments() []string {
	return []string{
		DocAccounts,
		DocBigMapActions,
		DocBigMapDiff,
		DocBigMapState,
		DocBlocks,
		DocContracts,
		DocGlobalConstants,
		DocMigrations,
		DocOperations,
		DocProtocol,
		DocScripts,
		DocTicketUpdates,
		DocTicketBalances,
		DocTickets,
		DocSmartRollups,
		DocStats,
	}
}

// AllModels -
func AllModels() []Model {
	return []Model{
		&protocol.Protocol{},
		&block.Block{},
		&account.Account{},
		&bigmapaction.BigMapAction{},
		&bigmapdiff.BigMapDiff{},
		&bigmapdiff.BigMapState{},
		&ticket.Ticket{},
		&ticket.TicketUpdate{},
		&ticket.Balance{},
		&operation.Operation{},
		&contract.GlobalConstant{},
		&contract.Script{},
		&contract.ScriptConstants{},
		&contract.Contract{},
		&migration.Migration{},
		&smartrollup.SmartRollup{},
		&stats.Stats{},
	}
}

// ManyToMany -
func ManyToMany() []interface{} {
	return []interface{}{
		&contract.ScriptConstants{},
	}
}
