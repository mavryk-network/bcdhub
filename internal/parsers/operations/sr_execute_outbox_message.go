package operations

import (
	"context"
	"encoding/hex"

	"github.com/mavryk-network/bcdhub/internal/bcd/encoding"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
	"github.com/pkg/errors"
)

// SrExecuteOutboxMessage -
type SrExecuteOutboxMessage struct {
	*ParseParams
}

// NewSrExecuteOutboxMessage -
func NewSrExecuteOutboxMessage(params *ParseParams) SrExecuteOutboxMessage {
	return SrExecuteOutboxMessage{params}
}

// Parse -
func (p SrExecuteOutboxMessage) Parse(ctx context.Context, data noderpc.Operation, store parsers.Store) error {
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
		Payload:      make([]byte, 0),
	}
	if data.Rollup != nil {
		operation.Destination = account.Account{
			Address:         *data.Rollup,
			Type:            types.NewAccountType(*data.Rollup),
			Level:           p.head.Level,
			OperationsCount: 1,
			LastAction:      p.head.Timestamp,
		}
		store.AddAccounts(operation.Destination)
	}
	p.fillInternal(&operation)
	operation.SetBurned(*p.protocol.Constants)
	parseOperationResult(data, &operation, store)
	p.stackTrace.Add(operation)

	if operation.IsApplied() {
		commitment, err := encoding.DecodeBase58(data.CementedCommitment)
		if err != nil {
			return errors.Wrap(err, "cemented commitment decoding")
		}
		operation.Payload = append(operation.Payload, commitment...)

		proof, err := hex.DecodeString(data.OutputProof)
		if err != nil {
			return errors.Wrap(err, "outbox proof decoding")
		}
		operation.Payload = append(operation.Payload, proof...)
	}

	store.AddOperations(&operation)
	store.AddAccounts(operation.Source)

	return nil
}

func (p SrExecuteOutboxMessage) fillInternal(tx *operation.Operation) {
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
