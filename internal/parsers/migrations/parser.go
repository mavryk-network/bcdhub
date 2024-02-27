package migrations

import (
	"context"
	"time"

	"github.com/mavryk-network/bcdhub/internal/models"
	modelsContract "github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/protocol"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
)

// Parser -
type Parser interface {
	Parse(ctx context.Context, script noderpc.Script, old *modelsContract.Contract, previous, next protocol.Protocol, timestamp time.Time, tx models.Transaction) error
	IsMigratable(address string) bool
}
