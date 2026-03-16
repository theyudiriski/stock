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
	"time"
)

func NewUpsertPriceIntraday(date, symbols string) *upsertPriceIntraday {
	ctx, cancel := context.WithCancel(context.Background())

	logger.Init()
	log := logger.Default

	stockbit := stockbit.NewStockbit(log, httpclient.New(service.ServiceNameStockbit))

	db, err := postgres.NewClient(config.LoadDatabase(), false)
	if err != nil {
		panic(err)
	}
	emittenStore := postgres.NewEmittenStore(db)
	priceFeedStore := postgres.NewPriceFeedStore(db)

	return &upsertPriceIntraday{
		base: base{
			logger:       log,
			emittenStore: emittenStore,
			symbols:      parseSymbols(symbols),
		},
		stockbit:       stockbit,
		priceFeedStore: priceFeedStore,
		ctx:            ctx,
		cancel:         cancel,

		date: date,
	}
}

type upsertPriceIntraday struct {
	base
	stockbit       service.Stockbit
	priceFeedStore service.PriceFeedStore

	ctx    context.Context
	cancel context.CancelFunc

	date string
}

func (u *upsertPriceIntraday) Run() (err error) {
	start := time.Now()
	ctx := u.ctx

	emittens, err := u.getEmittens(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	const maxConcurrency = 50
	sem := make(chan struct{}, maxConcurrency)

	errs := make(chan error, len(emittens))

	for i := range emittens {
		wg.Add(1)
		sem <- struct{}{}
		go func(emitten string) {
			defer wg.Done()
			defer func() { <-sem }()

			results, err := u.stockbit.GetPriceIntraday(ctx, emitten, u.date)
			if err != nil {
				errs <- fmt.Errorf("failed to get price intraday %s: %w", emitten, err)
				return
			}

			err = u.priceFeedStore.UpsertPriceIntraday(ctx, emitten, results)
			if err != nil {
				errs <- fmt.Errorf("failed to upsert price intraday %s: %w", emitten, err)
				return
			}

			u.logger.Info("successfully upserted price intraday", "symbol", emitten, "date", u.date)
		}(emittens[i])
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("failed to upsert price feed: %w", errors.Join(<-errs))
	}

	u.logger.Info("successfully upserted price intraday", "duration", time.Since(start), "date", u.date)
	return nil
}

func (u *upsertPriceIntraday) Stop() error {
	u.cancel()
	return nil
}
