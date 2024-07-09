package operations

import (
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
)

// GlobalConstant -
type GlobalConstant struct{}

// NewGlobalConstant -
func NewGlobalConstant() GlobalConstant {
	return GlobalConstant{}
}

// Parse -
func (p GlobalConstant) Parse(data noderpc.Operation, operation operation.Operation) *contract.GlobalConstant {
	gc := &contract.GlobalConstant{
		Timestamp: operation.Timestamp,
		Level:     operation.Level,
		Value:     data.Value,
	}

	if data.Metadata != nil && data.Metadata.OperationResult != nil {
		gc.Address = data.Metadata.OperationResult.GlobalAddress
	}

	return gc
}
