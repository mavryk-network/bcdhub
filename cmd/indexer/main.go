package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-io/workerpool"
	"github.com/grafana/pyroscope-go"
	"github.com/mavryk-network/bcdhub/cmd/indexer/indexer"
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/helpers"
	"github.com/mavryk-network/bcdhub/internal/logger"
	"github.com/mavryk-network/bcdhub/internal/profiler"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Err(err).Msg("loading config")
		return
	}

	logger.New(cfg.LogLevel)

	if cfg.Indexer.SentryEnabled {
		helpers.InitSentry(helpers.SentryConfig{
			DSN:   cfg.Sentry.URI,
			Debug: cfg.Sentry.Debug,
			Env:   cfg.Sentry.Environment,
			Tags: map[string]string{
				"project": cfg.Indexer.ProjectName,
			},
		})
		defer helpers.CatchPanicSentry()
	}

	var prof *pyroscope.Profiler
	if cfg.Profiler != nil {
		prof, err = profiler.New(cfg.Profiler.Server, "indexer")
		if err != nil {
			panic(err)
		}
	}

	g := workerpool.NewGroup()
	ctx, cancel := context.WithCancel(context.Background())

	indexers, err := indexer.CreateIndexers(ctx, cfg, g)
	if err != nil {
		cancel()
		log.Err(err).Msg("indexers creation")
		helpers.CatchErrorSentry(err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	cancel()

	g.Wait()
	for i := range indexers {
		if err := indexers[i].Close(); err != nil {
			panic(err)
		}
	}

	if prof != nil {
		if err := prof.Stop(); err != nil {
			panic(err)
		}
	}
	log.Info().Msg("stopped")
}
