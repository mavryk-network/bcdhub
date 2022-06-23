package search

import (
	"context"

	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/models/contract"
)

// SameContracts -
type SameContracts struct {
	Count     int64
	Contracts []Contract
}

// Searcher -
type Searcher interface {
	ByText(text string, offset int64, fields []string, filters map[string]interface{}, group bool) (Result, error)
	Save(ctx context.Context, items []Data) error
	CreateIndexes() error
	Rollback(network string, level int64) error
	BigMapDiffs(args BigMapDiffSearchArgs) ([]BigMapDiffResult, error)
	SameContracts(contract contract.Contract, network string, offset, size int64) (SameContracts, error)
}

// Data -
type Data interface {
	GetID() string
	GetIndex() string
}

// Constraint -
type Constraint[M models.Constraint] interface {
	BigMapDiff | Contract | Token | Metadata | Operation

	Data
}