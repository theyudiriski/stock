package cronjob

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"stock/config"
	"stock/internal/httpclient"
	"stock/internal/logger"
	"stock/internal/postgres"
	. "stock/internal/service"
	"stock/internal/stockbit"
	"sync"
)

func NewUpsertShareholderChart(symbols string) *upsertShareholderChart {
	ctx, cancel := context.WithCancel(context.Background())

	logger.Init()
	log := logger.Default

	stockbit := stockbit.NewStockbit(log, httpclient.New())

	db, err := postgres.NewClient(config.LoadDatabase(), false)
	if err != nil {
		panic(err)
	}

	emittenStore := postgres.NewEmittenStore(db)
	shareholderStore := postgres.NewShareholderStore(log, db)

	return &upsertShareholderChart{
		logger: log,

		stockbit:         stockbit,
		emittenStore:     emittenStore,
		shareholderStore: shareholderStore,

		ctx:    ctx,
		cancel: cancel,

		symbols: parseSymbols(symbols),
	}
}

type upsertShareholderChart struct {
	logger *slog.Logger

	stockbit         Stockbit
	emittenStore     EmittenStore
	shareholderStore ShareholderStore

	ctx    context.Context
	cancel context.CancelFunc

	symbols []string
}

func (u *upsertShareholderChart) Run() (err error) {
	ctx := u.ctx

	emittens := u.symbols
	if len(emittens) < 1 {
		emittens, err = u.emittenStore.GetEmittens(ctx)
		if err != nil {
			u.logger.Error("failed to get emittens", "error", err)
			return err
		}
	}

	var wg sync.WaitGroup
	const maxConcurrency = 10
	sem := make(chan struct{}, maxConcurrency)

	errs := make(chan error, len(emittens))

	for i := range emittens {
		wg.Add(1)
		sem <- struct{}{}
		go func(emitten string) {
			defer wg.Done()
			defer func() { <-sem }()

			shareholderData, err := stockbit.Retryable(func() (*ShareholderChartData, error) {
				return u.stockbit.GetShareholders(ctx, emitten, nil)
			}).WithRetry()
			if err != nil {
				errs <- fmt.Errorf("failed to get shareholder chart %s: %w", emitten, err)
				return
			}

			history, err := u.shareholderStore.GetShareholderChartHistory(ctx, emitten)
			if err != nil {
				errs <- fmt.Errorf("failed to get shareholder chart history %s: %w", emitten, err)
				return
			}
			if history != nil && *history == shareholderData.LastUpdate {
				u.logger.Info("shareholder chart history is the same as the last updated date", "symbol", emitten)
				return
			}

			err = u.shareholderStore.UpsertShareholderChart(ctx, emitten, shareholderData)
			if err != nil {
				errs <- fmt.Errorf("failed to upsert shareholder chart %s: %w", emitten, err)
				return
			}

			u.logger.Info("successfully upserted shareholder chart", "symbol", emitten)
		}(emittens[i])
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("failed to upsert subsidiary companies: %w", errors.Join(<-errs))
	}

	u.logger.Info("successfully upserted subsidiary companies")
	return nil
}

func (u *upsertShareholderChart) Stop() error {
	u.cancel()
	return nil
}
