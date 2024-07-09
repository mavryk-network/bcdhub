package core

import (
	"context"
	"sync"

	"github.com/mavryk-network/bcdhub/internal/models"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapaction"
	"github.com/mavryk-network/bcdhub/internal/models/bigmapdiff"
	"github.com/mavryk-network/bcdhub/internal/models/block"
	"github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/operation"
	"github.com/mavryk-network/bcdhub/internal/models/ticket"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

func createTable(ctx context.Context, db bun.IDB, model models.Model) error {
	if model == nil {
		return nil
	}

	query := db.
		NewCreateTable().
		Model(model).
		IfNotExists()

	_, err := query.Exec(ctx)
	return err
}

func createTables(ctx context.Context, db *bun.DB) error {
	for _, model := range models.AllModels() {
		if err := createTable(ctx, db, model); err != nil {
			return err
		}
	}
	return createHypertables(ctx, db)
}

func createHypertables(ctx context.Context, db *bun.DB) error {
	for _, model := range []models.Model{
		&block.Block{},
		&bigmapdiff.BigMapDiff{},
		&bigmapaction.BigMapAction{},
		&contract.Contract{},
		&migration.Migration{},
		&operation.Operation{},
		&ticket.TicketUpdate{},
	} {
		if _, err := db.ExecContext(ctx,
			`SELECT public.create_hypertable(?, 'timestamp', chunk_time_interval => INTERVAL '1 month', if_not_exists => TRUE);`,
			model.TableName(),
		); err != nil {
			return err
		}
	}

	return nil
}

var createExtensionOnce sync.Once
var wg sync.WaitGroup

func createSchema(ctx context.Context, db *bun.DB, schemaName string) error {
	createExtensionOnce.Do(func() {
		defer wg.Done()
		wg.Add(1)
		if _, err := db.NewRaw("set search_path = 'public'").Exec(ctx); err != nil {
			log.Err(err).Msg("set search path to public")
			return
		}
		if _, err := db.NewRaw("CREATE EXTENSION IF NOT EXISTS timescaledb;").Exec(ctx); err != nil {
			log.Err(err).Msg("create timescale extension")
		}
	})

	wg.Wait()

	schema := bun.Ident(schemaName)
	if _, err := db.NewRaw("create schema if not exists ?", schema).Exec(ctx); err != nil {
		return err
	}
	if _, err := db.NewRaw("set search_path = ?", schema).Exec(ctx); err != nil {
		return err
	}
	return nil
}

// Drop - drops full database
func (p *Postgres) Drop(ctx context.Context) error {
	for _, table := range models.ManyToMany() {
		if _, err := p.DB.NewDropTable().Model(table).IfExists().Cascade().Exec(ctx); err != nil {
			return err
		}
	}

	for _, table := range models.AllModels() {
		if _, err := p.DB.NewDropTable().Model(table).IfExists().Cascade().Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

const tableExistsQuery = `SELECT EXISTS(
    SELECT * 
    FROM information_schema.tables 
    WHERE 
      table_schema = ? AND 
      table_name = ?
) as flag;`

type existsResponse struct {
	Flag bool `bun:"flag"`
}

// TablesExist - returns true if all tables exist otherwise false
func (p *Postgres) TablesExist(ctx context.Context) bool {
	for _, table := range models.AllDocuments() {
		var exists existsResponse
		err := p.DB.NewRaw(tableExistsQuery, p.schema, table).Scan(ctx, &exists)
		if !exists.Flag || err != nil {
			return false
		}
	}
	return true
}
