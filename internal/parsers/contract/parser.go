package contract

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/parsers"
)

// Parser -
type Parser interface {
	Parse(ctx context.Context, operation *operation.Operation, store parsers.Store) error
}
