package contract

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/bcd"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/postgres/core"
	"github.com/pkg/errors"
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
func (storage *Storage) Get(ctx context.Context, address string) (response contract.Contract, err error) {
	var accountID int64
	if err = storage.DB.NewSelect().
		Model((*account.Account)(nil)).
		Column("id").
		Where("address = ?", address).
		Scan(ctx, &accountID); err != nil {
		return
	}

	err = storage.DB.NewSelect().
		Model(&response).
		Where("contract.account_id = ?", accountID).
		Relation("Account").Relation("Manager").
		Relation("Delegate").Relation("Alpha").
		Relation("Babylon").Relation("Jakarta").
		Scan(ctx)
	return
}

// ByHash -
func (storage *Storage) ByHash(ctx context.Context, hash string) (result contract.Script, err error) {
	err = storage.DB.NewSelect().Model(&result).Where("hash = ?", hash).Limit(1).Scan(ctx)
	return
}

// Script -
func (storage *Storage) Script(ctx context.Context, address string, symLink string) (contract.Script, error) {
	var accountID int64
	if err := storage.DB.
		NewSelect().
		Model((*account.Account)(nil)).
		Column("id").
		Where("address = ?", address).
		Scan(ctx, &accountID); err != nil {
		return contract.Script{}, err
	}

	var c contract.Contract
	query := storage.DB.NewSelect().Model(&c).Where("account_id = ?", accountID)
	switch symLink {
	case bcd.SymLinkAlpha:
		err := query.Relation("Alpha").Scan(ctx)
		return c.Alpha, err
	case bcd.SymLinkBabylon:
		err := query.Relation("Babylon").Scan(ctx)
		return c.Babylon, err
	case bcd.SymLinkJakarta:
		err := query.Relation("Jakarta").Scan(ctx)
		return c.Jakarta, err
	}
	return c.Alpha, errors.Errorf("unknown protocol symbolic link: %s", symLink)
}

// Code -
func (storage *Storage) Code(ctx context.Context, id int64) ([]byte, error) {
	var data []byte
	err := storage.DB.NewSelect().Model((*contract.Script)(nil)).Where("id = ?", id).Column("code").Scan(ctx, &data)
	return data, err
}

// Parameter -
func (storage *Storage) Parameter(ctx context.Context, id int64) ([]byte, error) {
	var data []byte
	err := storage.DB.NewSelect().
		Model((*contract.Script)(nil)).
		Where("id = ?", id).
		Column("parameter").
		Scan(ctx, &data)
	return data, err
}

// Storage -
func (storage *Storage) Storage(ctx context.Context, id int64) ([]byte, error) {
	var data []byte
	err := storage.DB.NewSelect().
		Model((*contract.Script)(nil)).
		Where("id = ?", id).
		Column("storage").
		Scan(ctx, &data)
	return data, err
}

// Views -
func (storage *Storage) Views(ctx context.Context, id int64) ([]byte, error) {
	var data []byte
	err := storage.DB.NewSelect().
		Model((*contract.Script)(nil)).
		Where("id = ?", id).
		Column("views").
		Scan(ctx, &data)
	return data, err
}

// ScriptPart -
func (storage *Storage) ScriptPart(ctx context.Context, address string, symLink, part string) ([]byte, error) {
	var accountID int64
	if err := storage.DB.NewSelect().
		Model((*account.Account)(nil)).
		Column("id").
		Where("address = ?", address).
		Scan(ctx, &accountID); err != nil {
		return nil, err
	}

	var scriptId int64
	scriptIdQuery := storage.DB.NewSelect().
		Model((*contract.Contract)(nil)).
		Where("account_id = ?", accountID)

	switch symLink {
	case bcd.SymLinkAlpha:
		scriptIdQuery = scriptIdQuery.Column("alpha_id")
	case bcd.SymLinkBabylon:
		scriptIdQuery = scriptIdQuery.Column("babylon_id")
	case bcd.SymLinkJakarta:
		scriptIdQuery = scriptIdQuery.Column("jakarta_id")
	default:
		return nil, errors.Errorf("unknown protocol symbolic link: %s", symLink)
	}

	if err := scriptIdQuery.Scan(ctx, &scriptId); err != nil {
		return nil, err
	}

	partQuery := storage.DB.NewSelect().
		Model((*contract.Script)(nil)).
		Where("id = ?", scriptId)

	switch part {
	case consts.PARAMETER:
		partQuery.Column("parameter")
	case consts.CODE:
		partQuery.Column("code")
	case consts.STORAGE:
		partQuery.Column("storage")
	case consts.VIEWS:
		partQuery.Column("views")
	default:
		return nil, errors.Errorf("unknown script part name: %s", part)
	}

	var data []byte
	err := partQuery.Scan(ctx, &data)
	return data, err
}

// FindOne -
func (storage *Storage) FindOne(ctx context.Context, tags types.Tags) (result contract.Contract, err error) {
	err = storage.DB.NewSelect().Model(&result).
		Where("tags&? > 0", tags).
		ColumnExpr("contract.id, contract.account_id, contract.timestamp, contract.level").
		ColumnExpr("account.address as account__address, account.operations_count as account__operations_count, account.last_action as account__last_action").
		Join(`LEFT JOIN "accounts" AS "account" ON "account"."id" = "contract"."account_id"`).
		Limit(1).
		Scan(ctx)
	return
}

func (storage *Storage) AllExceptDelegators(ctx context.Context) (contracts []contract.Contract, err error) {
	err = storage.DB.NewSelect().Model(&contracts).
		Relation("Account").
		Where("tags & 4 = 0"). // except delegator contracts
		Scan(ctx)
	return
}
