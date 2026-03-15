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
	"stock/internal/service"
	"stock/internal/stockbit"
	"sync"
)

func NewUpsertCorporateAction(symbols string) *upsertCorporateAction {
	ctx, cancel := context.WithCancel(context.Background())

	logger.Init()
	log := logger.Default

	stockbit := stockbit.NewStockbit(log, httpclient.New(service.ServiceNameStockbit))

	db, err := postgres.NewClient(config.LoadDatabase(), false)
	if err != nil {
		panic(err)
	}

	emittenStore := postgres.NewEmittenStore(db)
	corpactionStore := postgres.NewCorpactionStore(db)

	return &upsertCorporateAction{
		logger:          log,
		stockbit:        stockbit,
		emittenStore:    emittenStore,
		corpactionStore: corpactionStore,
		ctx:             ctx,
		cancel:          cancel,

		symbols: parseSymbols(symbols),
	}
}

type upsertCorporateAction struct {
	logger *slog.Logger

	stockbit        service.Stockbit
	emittenStore    service.EmittenStore
	corpactionStore service.CorpactionStore

	ctx    context.Context
	cancel context.CancelFunc

	symbols []string
}

func (u *upsertCorporateAction) Run() (err error) {
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
	const maxConcurrency = 4
	sem := make(chan struct{}, maxConcurrency)

	errs := make(chan error, len(emittens))

	for i := range emittens {
		wg.Add(1)
		sem <- struct{}{}
		go func(emitten string) {
			defer wg.Done()
			defer func() { <-sem }()

			// corpaction, err := retry.DoWithData(
			// 	func() (*service.Corpaction, error) {
			// 		return u.stockbit.GetCorporateActions(ctx, emitten)
			// 	},
			// 	retry.Attempts(5),
			// 	retry.Delay(500*time.Millisecond),
			// 	retry.OnRetry(func(n uint, err error) {
			// 		u.logger.Warn("failed to get corporate actions, retrying...", "symbol", emitten, "attempt", n+1, "error", err)
			// 	}),
			// )
			// if err != nil {
			// 	errs <- fmt.Errorf("failed to get corporate actions %s: %w", emitten, err)
			// 	return
			// }

			// err = u.corpactionStore.UpsertCorpactions(ctx, emitten, corpaction)
			// if err != nil {
			// 	errs <- fmt.Errorf("failed to upsert corporate actions %s: %w", emitten, err)
			// 	return
			// }

			u.logger.Info("successfully upserted corporate actions", "symbol", emitten)
		}(emittens[i])
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("failed to upsert corporate actions: %w", errors.Join(<-errs))
	}

	u.logger.Info("successfully upserted corporate actions")
	return nil
}

func (u *upsertCorporateAction) Stop() error {
	u.cancel()
	return nil
}
