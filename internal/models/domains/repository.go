package domains

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models/contract"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/domains/mock.go -package=domains -typed
type Repository interface {
	Same(ctx context.Context, network string, c contract.Contract, limit, offset int, availiableNetworks ...string) ([]Same, error)
	SameCount(ctx context.Context, c contract.Contract, availiableNetworks ...string) (int, error)
}
