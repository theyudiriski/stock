package cronjob

import (
	"context"
	"errors"
	"fmt"
	"stock/config"
	"stock/internal/httpclient"
	"stock/internal/logger"
	"stock/internal/postgres"
	"stock/internal/service"
	"stock/internal/stockbit"
	"sync"
)

func NewUpsertEmittenCompanySubsidiary(symbols string) *upsertEmittenCompanySubsidiary {
	ctx, cancel := context.WithCancel(context.Background())

	logger.Init()
	log := logger.Default

	stockbit := stockbit.NewStockbit(log, httpclient.New(service.ServiceNameStockbit))

	db, err := postgres.NewClient(config.LoadDatabase(), false)
	if err != nil {
		panic(err)
	}

	emittenStore := postgres.NewEmittenStore(db)
	subsidiaryStore := postgres.NewSubsidiaryStore(db, log)

	return &upsertEmittenCompanySubsidiary{
		base: base{
			logger:       log,
			emittenStore: emittenStore,
			symbols:      parseSymbols(symbols),
		},
		stockbit:        stockbit,
		subsidiaryStore: subsidiaryStore,

		ctx:    ctx,
		cancel: cancel,
	}
}

type upsertEmittenCompanySubsidiary struct {
	base
	stockbit        service.Stockbit
	subsidiaryStore service.SubsidiaryStore

	ctx    context.Context
	cancel context.CancelFunc
}

func (u *upsertEmittenCompanySubsidiary) Run() (err error) {
	ctx := u.ctx

	emittens, err := u.getEmittens(ctx)
	if err != nil {
		return err
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

			subsidiaryData, err := stockbit.Retryable(func() (*service.SubsidiaryData, error) {
				return u.stockbit.GetSubsidiaryCompanies(ctx, emitten)
			}).WithRetry()
			if err != nil {
				errs <- fmt.Errorf("failed to get subsidiary companies %s: %w", emitten, err)
				return
			}

			if subsidiaryData != nil {
				history, err := u.subsidiaryStore.GetSubsidiaryCompanyHistory(ctx, emitten)
				if err != nil {
					errs <- fmt.Errorf("failed to get subsidiary company history %s: %w", emitten, err)
					return
				}
				if history != nil && *history == subsidiaryData.LastUpdatedPeriod {
					u.logger.Info("subsidiary company history is the same as the last updated period", "symbol", emitten)
					return
				}

				err = u.subsidiaryStore.UpsertSubsidiaryCompanies(ctx, emitten, subsidiaryData)
				if err != nil {
					errs <- fmt.Errorf("failed to upsert subsidiary companies %s: %w", emitten, err)
					return
				}

				u.logger.Info("successfully upserted subsidiary companies", "symbol", emitten)
			} else {
				u.logger.Info("no subsidiary companies found", "symbol", emitten)
			}
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

func (u *upsertEmittenCompanySubsidiary) Stop() error {
	u.cancel()
	return nil
}
