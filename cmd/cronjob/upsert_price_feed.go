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

func NewUpsertPriceFeed(fromDate, toDate, symbols string) *upsertPriceFeed {
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

	return &upsertPriceFeed{
		base: base{
			logger:       log,
			emittenStore: emittenStore,
			symbols:      parseSymbols(symbols),
		},
		stockbit:       stockbit,
		priceFeedStore: priceFeedStore,
		ctx:            ctx,
		cancel:         cancel,

		fromDate: fromDate,
		toDate:   toDate,
	}
}

type upsertPriceFeed struct {
	base
	stockbit       service.Stockbit
	priceFeedStore service.PriceFeedStore

	ctx    context.Context
	cancel context.CancelFunc

	fromDate string
	toDate   string
}

func (u *upsertPriceFeed) Run() (err error) {
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

			priceFeed, err := u.stockbit.GetPriceFeed(ctx, emitten, u.fromDate, u.toDate)
			if err != nil {
				errs <- fmt.Errorf("failed to get price feed %s: %w", emitten, err)
				return
			}

			err = u.priceFeedStore.UpsertPriceFeed(ctx, emitten, priceFeed)
			if err != nil {
				errs <- fmt.Errorf("failed to upsert price feed %s: %w", emitten, err)
				return
			}

			u.logger.Info("successfully upserted price feed", "symbol", emitten)
		}(emittens[i])
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("failed to upsert price feed: %w", errors.Join(<-errs))
	}

	u.logger.Info("successfully upserted price feed", "duration", time.Since(start), "fromDate", u.fromDate, "toDate", u.toDate)
	return nil
}

func (u *upsertPriceFeed) Stop() error {
	u.cancel()
	return nil
}
