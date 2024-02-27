package migrations

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/protocol"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
	"github.com/mavryk-network/bcdhub/internal/parsers/contract"
)

// VestingParser -
type VestingParser struct {
	parser   contract.Parser
	protocol protocol.Protocol
}

// NewVestingParser -
func NewVestingParser(ctx *config.Context, contractParser contract.Parser, proto protocol.Protocol) *VestingParser {
	return &VestingParser{
		parser:   contractParser,
		protocol: proto,
	}
}

// Parse -
func (p *VestingParser) Parse(ctx context.Context, data noderpc.ContractData, head noderpc.Header, address string, store parsers.Store) error {
	vestingOperation := &operation.Operation{
		ProtocolID: p.protocol.ID,
		Status:     types.OperationStatusApplied,
		Kind:       types.OperationKindOrigination,
		Amount:     data.Balance,
		Counter:    data.Counter,
		Source: account.Account{
			Address:    data.Manager,
			Type:       types.NewAccountType(data.Manager),
			Level:      head.Level,
			LastAction: head.Timestamp,
		},
		Destination: account.Account{
			Address:         address,
			Type:            types.NewAccountType(address),
			Level:           head.Level,
			LastAction:      head.Timestamp,
			MigrationsCount: 1,
		},
		Delegate: account.Account{
			Address:    data.Delegate.Value,
			Type:       types.NewAccountType(data.Delegate.Value),
			Level:      head.Level,
			LastAction: head.Timestamp,
		},
		Level:     head.Level,
		Timestamp: head.Timestamp,
		Script:    data.RawScript,
	}
	if err := p.parser.Parse(ctx, vestingOperation, store); err != nil {
		return err
	}

	contracts := store.ListContracts()
	for i := range contracts {
		if contracts[i].Account.Address == address {
			store.AddMigrations(&migration.Migration{
				Level:      head.Level,
				ProtocolID: p.protocol.ID,
				Timestamp:  head.Timestamp,
				Kind:       types.MigrationKindBootstrap,
				Contract:   *contracts[i],
			})
			store.AddAccounts(
				vestingOperation.Source,
				vestingOperation.Destination,
				vestingOperation.Delegate,
			)
			break
		}
	}

	return nil
}
