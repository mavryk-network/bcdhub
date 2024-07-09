package operations

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
)

// SrOriginate -
type SrOriginate struct {
	*ParseParams
}

// NewSrOriginate -
func NewSrOriginate(params *ParseParams) SrOriginate {
	return SrOriginate{params}
}

// Parse -
func (p SrOriginate) Parse(ctx context.Context, data noderpc.Operation, store parsers.Store) error {
	source := account.Account{
		Address:         data.Source,
		Type:            types.NewAccountType(data.Source),
		Level:           p.head.Level,
		OperationsCount: 1,
		LastAction:      p.head.Timestamp,
	}

	operation := operation.Operation{
		Source:       source,
		Initiator:    source,
		Fee:          data.Fee,
		Counter:      data.Counter,
		StorageLimit: data.StorageLimit,
		GasLimit:     data.GasLimit,
		Hash:         p.hash,
		ProtocolID:   p.protocol.ID,
		Level:        p.head.Level,
		Timestamp:    p.head.Timestamp,
		Kind:         types.NewOperationKind(data.Kind),
		ContentIndex: p.contentIdx,
	}

	p.fillInternal(&operation)
	operation.SetBurned(*p.protocol.Constants)
	parseOperationResult(data, &operation, store)
	p.stackTrace.Add(operation)

	store.AddOperations(&operation)

	if operation.IsApplied() {
		smartRollup, err := NewSmartRolupParser().Parse(data, operation)
		if err != nil {
			return err
		}
		store.AddSmartRollups(&smartRollup)
		store.AddAccounts(smartRollup.Address)
		operation.Destination = smartRollup.Address
	}

	store.AddAccounts(operation.Source)

	return nil
}

func (p SrOriginate) fillInternal(tx *operation.Operation) {
	if p.main == nil {
		p.main = tx
		return
	}

	tx.Counter = p.main.Counter
	tx.Hash = p.main.Hash
	tx.Level = p.main.Level
	tx.Timestamp = p.main.Timestamp
	tx.Internal = true
	tx.Initiator = p.main.Source
}
