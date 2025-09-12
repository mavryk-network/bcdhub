package storage

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
	"github.com/mavryk-network/bcdhub/internal/parsers"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Parser -
type Parser interface {
	ParseTransaction(ctx context.Context, content noderpc.Operation, operation *operation.Operation, store parsers.Store) error
	ParseOrigination(ctx context.Context, content noderpc.Operation, operation *operation.Operation, store parsers.Store) error
}
