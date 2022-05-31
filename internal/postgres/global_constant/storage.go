package global_constant

import (
	"github.com/baking-bad/bcdhub/internal/models/contract"
	"github.com/baking-bad/bcdhub/internal/postgres/core"
	"github.com/go-pg/pg/v10"
)

// Storage -
type Storage struct {
	*core.Postgres
}

// NewStorage -
func NewStorage(pg *core.Postgres) *Storage {
	return &Storage{pg}
}

// Get -
func (storage *Storage) Get(address string) (response contract.GlobalConstant, err error) {
	query := storage.DB.Model(&response)
	core.Address(address)(query)
	err = query.First()
	return
}

// All -
func (storage *Storage) All(addresses ...string) (response []contract.GlobalConstant, err error) {
	if len(addresses) == 0 {
		return
	}

	err = storage.DB.Model(new(contract.GlobalConstant)).Where("address IN (?)", pg.In(addresses)).Select(&response)
	return
}