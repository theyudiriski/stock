package cronjob

import (
	"context"
	"errors"
	"fmt"
	"stock/config"
	"stock/internal/cache"
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

	holidayStore := postgres.NewHolidayStore(db)
	redisClient := cache.NewRedisClient(config.LoadRedis())
	holidayCache := cache.NewHolidayCache(redisClient, holidayStore)

	return &upsertPriceIntraday{
		base: base{
			logger:       log,
			emittenStore: emittenStore,
			symbols:      parseSymbols(symbols),
		},
		stockbit:       stockbit,
		priceFeedStore: priceFeedStore,
		holidayCache:   holidayCache,

		ctx:    ctx,
		cancel: cancel,

		date: date,
	}
}

type upsertPriceIntraday struct {
	base
	stockbit       service.Stockbit
	priceFeedStore service.PriceFeedStore
	holidayCache   service.HolidayStore // store layer contains cache and postgres

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

	date, err := time.Parse(time.DateOnly, u.date)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to parse date", "date", u.date, "error", err)
		return err
	}

	holidays, err := u.holidayCache.GetHolidaySet(ctx, date, date)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get holidays", "date", date, "error", err)
		return err
	}

	if holidays[date.Format(time.DateOnly)] {
		u.logger.InfoContext(ctx, "skipping price intraday because it is a holiday", "date", date)
		return nil
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
