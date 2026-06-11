package migrations

import (
	"context"
	"time"

	"github.com/mavryk-network/nexushub/internal/models"
	modelsContract "github.com/mavryk-network/nexushub/internal/models/contract"
	"github.com/mavryk-network/nexushub/internal/models/protocol"
	"github.com/mavryk-network/nexushub/internal/noderpc"
)

// Parser -
type Parser interface {
	Parse(ctx context.Context, script noderpc.Script, old *modelsContract.Contract, previous, next protocol.Protocol, timestamp time.Time, tx models.Transaction) error
	IsMigratable(address string) bool
}
