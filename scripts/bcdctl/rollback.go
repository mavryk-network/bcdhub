package main

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/postgres"
	"github.com/mavryk-network/bcdhub/internal/rollback"
	"github.com/rs/zerolog/log"
)

type rollbackCommand struct {
	Level   int64  `description:"Level to rollback" long:"level"   short:"l"`
	Network string `description:"Network"           long:"network" short:"n"`
}

var rollbackCmd rollbackCommand

// Execute
func (x *rollbackCommand) Execute(_ []string) error {
	network := types.NewNetwork(x.Network)
	ctx, err := ctxs.Get(network)
	if err != nil {
		panic(err)
	}

	state, err := ctx.Blocks.Last(context.Background())
	if err != nil {
		panic(err)
	}

	log.Warn().Msgf("Do you want to rollback '%s' from %d to %d? (yes - continue. no - cancel)", network.String(), state.Level, x.Level)
	if !yes() {
		log.Info().Msg("Cancelled")
		return nil
	}

	if err := ctx.Storage.InitDatabase(context.Background()); err != nil {
		return err
	}

	saver, err := postgres.NewRollback(ctx.StorageDB.DB)
	if err != nil {
		return err
	}
	manager := rollback.NewManager(ctx.Storage, ctx.Blocks, saver, ctx.Stats)
	if err = manager.Rollback(context.Background(), network, state, x.Level); err != nil {
		return err
	}
	log.Info().Msg("Done")

	return nil
}
