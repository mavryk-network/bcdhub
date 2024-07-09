package contract

import (
	"context"
	"time"

	"github.com/mavryk-network/bcdhub/internal/models/types"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/contract/mock.go -package=contract -typed
type Repository interface {
	Get(ctx context.Context, address string) (Contract, error)
	Script(ctx context.Context, address string, symLink string) (Script, error)

	// ScriptPart - returns part of script type. Part can be `storage`, `parameter` or `code`.
	ScriptPart(ctx context.Context, address string, symLink, part string) ([]byte, error)
	FindOne(ctx context.Context, tags types.Tags) (Contract, error)
	AllExceptDelegators(ctx context.Context) ([]Contract, error)
}

//go:generate mockgen -source=$GOFILE -destination=../mock/contract/mock.go -package=contract -typed
type ScriptRepository interface {
	ByHash(ctx context.Context, hash string) (Script, error)
	Code(ctx context.Context, id int64) ([]byte, error)
	Parameter(ctx context.Context, id int64) ([]byte, error)
	Storage(ctx context.Context, id int64) ([]byte, error)
	Views(ctx context.Context, id int64) ([]byte, error)
}

//go:generate mockgen -source=$GOFILE -destination=../mock/contract/mock.go -package=contract -typed
type ConstantRepository interface {
	Get(ctx context.Context, address string) (GlobalConstant, error)
	All(ctx context.Context, addresses ...string) ([]GlobalConstant, error)
	List(ctx context.Context, size, offset int64, orderBy, sort string) ([]ListGlobalConstantItem, error)
	ForContract(ctx context.Context, address string, size, offset int64) ([]GlobalConstant, error)
	ContractList(ctx context.Context, address string, size, offset int64) ([]Contract, error)
}

// ListGlobalConstantItem -
type ListGlobalConstantItem struct {
	Timestamp  time.Time `bun:"timestamp"`
	Level      int64     `bun:"level"`
	Address    string    `bun:"address"`
	LinksCount uint64    `bun:"links_count"`
}
