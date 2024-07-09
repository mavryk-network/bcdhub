package migrations

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	contracts "github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/protocol"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
	"github.com/mavryk-network/bcdhub/internal/parsers/contract"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// ImplicitParser -
type ImplicitParser struct {
	ctx            *config.Context
	rpc            noderpc.INode
	contractParser contract.Parser
	protocol       protocol.Protocol
	contractsRepo  contracts.Repository
}

// NewImplicitParser -
func NewImplicitParser(ctx *config.Context,
	rpc noderpc.INode,
	contractParser contract.Parser,
	protocol protocol.Protocol,
	contractsRepo contracts.Repository) (*ImplicitParser, error) {
	return &ImplicitParser{ctx, rpc, contractParser, protocol, contractsRepo}, nil
}

// Parse -
func (p *ImplicitParser) Parse(ctx context.Context, metadata noderpc.Metadata, head noderpc.Header, store parsers.Store) error {
	if len(metadata.ImplicitOperationsResults) == 0 {
		return nil
	}

	for i := range metadata.ImplicitOperationsResults {
		switch metadata.ImplicitOperationsResults[i].Kind {
		case consts.Origination:
			if err := p.origination(ctx, metadata.ImplicitOperationsResults[i], head, store); err != nil {
				return err
			}
		case consts.Transaction:
			if err := p.transaction(metadata.ImplicitOperationsResults[i], head, store); err != nil {
				return err
			}
		}
	}
	return nil
}

// IsMigratable -
func (p *ImplicitParser) IsMigratable(address string) bool {
	return true
}

func (p *ImplicitParser) origination(ctx context.Context, implicit noderpc.ImplicitOperationsResult, head noderpc.Header, store parsers.Store) error {
	if len(implicit.OriginatedContracts) == 0 {
		return nil
	}
	if _, err := p.contractsRepo.Get(ctx, implicit.OriginatedContracts[0]); err == nil {
		return nil
	}
	origination := operation.Operation{
		ProtocolID: p.protocol.ID,
		Level:      head.Level,
		Timestamp:  head.Timestamp,
		Kind:       types.OperationKindOrigination,
		Destination: account.Account{
			Address:         implicit.OriginatedContracts[0],
			Type:            types.AccountTypeContract,
			Level:           head.Level,
			OperationsCount: 1,
			LastAction:      head.Timestamp,
			MigrationsCount: 1,
		},
		ConsumedGas:         implicit.ConsumedGas,
		PaidStorageSizeDiff: implicit.PaidStorageSizeDiff,
		StorageSize:         implicit.StorageSize,
		DeffatedStorage:     implicit.Storage,
	}

	script, err := p.rpc.GetRawScript(ctx, origination.Destination.Address, origination.Level)
	if err != nil {
		return err
	}
	origination.Script = script

	if err := p.contractParser.Parse(ctx, &origination, store); err != nil {
		return err
	}

	contracts := store.ListContracts()
	for i := range contracts {
		if contracts[i].Account.Address == implicit.OriginatedContracts[0] {
			store.AddMigrations(&migration.Migration{
				ProtocolID: p.protocol.ID,
				Level:      head.Level,
				Timestamp:  head.Timestamp,
				Kind:       types.MigrationKindBootstrap,
				Contract:   *contracts[i],
			})
			store.AddAccounts(origination.Destination)
			break
		}
	}

	log.Info().Msg("Implicit bootstrap migration found")
	return nil
}

func (p *ImplicitParser) transaction(implicit noderpc.ImplicitOperationsResult, head noderpc.Header, store parsers.Store) error {
	tx := operation.Operation{
		ProtocolID:      p.protocol.ID,
		Level:           head.Level,
		Timestamp:       head.Timestamp,
		Kind:            types.OperationKindTransaction,
		ConsumedGas:     implicit.ConsumedGas,
		StorageSize:     implicit.StorageSize,
		DeffatedStorage: implicit.Storage,
		Status:          types.OperationStatusApplied,
		Tags:            types.NewTags([]string{types.ImplicitOperationStringTag}),
		Counter:         head.Level,
	}

	for i := range implicit.BalanceUpdates {
		if implicit.BalanceUpdates[i].Kind == "contract" && implicit.BalanceUpdates[i].Origin == "subsidy" {
			tx.Destination = account.Account{
				Type:            types.NewAccountType(implicit.BalanceUpdates[i].Contract),
				Address:         implicit.BalanceUpdates[i].Contract,
				Level:           head.Level,
				OperationsCount: 1,
				LastAction:      head.Timestamp,
			}
			store.AddAccounts(tx.Destination)
			break
		}
	}

	if tx.Destination.Address == "" {
		return errors.Errorf("empty destination in implicit transaction at level %d", head.Level)
	}

	store.AddOperations(&tx)
	return nil
}
