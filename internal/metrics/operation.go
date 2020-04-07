package metrics

import (
	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/models"
)

// SetOperationAliases -
func (h *Handler) SetOperationAliases(aliases map[string]string, op *models.Operation) {
	op.SourceAlias = aliases[op.Source]
	op.DestinationAlias = aliases[op.Destination]

	if op.Delegate != "" {
		op.DelegateAlias = aliases[op.Delegate]
	}
}

// SetOperationBurned -
func (h *Handler) SetOperationBurned(op *models.Operation) {
	if op.Status != consts.Applied {
		return
	}

	if op.Result == nil {
		return
	}

	var burned int64

	if op.Result.PaidStorageSizeDiff != 0 {
		burned += op.Result.PaidStorageSizeDiff * 1000
	}

	if op.Result.AllocatedDestinationContract {
		burned += 257000
	}

	op.Burned = burned
}