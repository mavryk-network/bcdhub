package operations

import (
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/mavryk-network/bcdhub/internal/bcd/tezerrors"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
)

func parseOperationResult(data noderpc.Operation, tx *operation.Operation, store parsers.Store) {
	result := data.GetResult()
	if result == nil {
		return
	}

	tx.Status = types.NewOperationStatus(result.Status)

	if result.ConsumedMilligas != nil {
		tx.ConsumedGas = *result.ConsumedMilligas
	} else {
		tx.ConsumedGas = result.ConsumedGas * 100
	}

	if result.StorageSize != nil {
		tx.StorageSize = *result.StorageSize
	}
	if result.PaidStorageSizeDiff != nil {
		tx.PaidStorageSizeDiff = *result.PaidStorageSizeDiff
	}
	if len(result.Originated) > 0 {
		tx.Destination = account.Account{
			Address:         result.Originated[0],
			Type:            types.AccountTypeContract,
			Level:           tx.Level,
			OperationsCount: 1,
			LastAction:      tx.Timestamp,
		}
	}

	if len(result.OriginatedRollup) > 0 {
		tx.Destination = account.Account{
			Address:         result.OriginatedRollup,
			Type:            types.AccountTypeRollup,
			Level:           tx.Level,
			OperationsCount: 1,
			LastAction:      tx.Timestamp,
		}
	}

	tx.AllocatedDestinationContract = data.Kind == consts.Origination
	if result.AllocatedDestinationContract != nil {
		tx.AllocatedDestinationContract = *result.AllocatedDestinationContract
	}

	if errs, err := tezerrors.ParseArray(result.Errors); err == nil {
		tx.Errors = errs
	}

	if tx.IsApplied() {
		new(TicketUpdateParser).Parse(result, tx, store)
	}
}
