package operation

import (
	"time"

	"github.com/mavryk-network/bcdhub/internal/models/types"
)

// DAppStats -
type DAppStats struct {
	Users  int64 `json:"users"`
	Calls  int64 `json:"txs"`
	Volume int64 `json:"volume"`
}

// OPG -
type OPG struct {
	LastID       int64
	ContentIndex int64
	Counter      int64
	Level        int64
	TotalCost    int64
	Flow         int64
	Internals    int
	Hash         []byte
	Entrypoint   string
	Timestamp    time.Time
	Status       types.OperationStatus
	Kind         types.OperationKind
}

// ContractStats -
type ContractStats struct {
	Count      int64
	LastAction time.Time
}
