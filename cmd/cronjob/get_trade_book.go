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

func NewGetTradeBook(symbols string) *getTradeBook {
	ctx, cancel := context.WithCancel(context.Background())

	logger.Init()
	log := logger.Default

	db, err := postgres.NewClient(config.LoadDatabase(), false)
	if err != nil {
		panic(err)
	}
	emittenStore := postgres.NewEmittenStore(db)
	stockbit := stockbit.NewStockbit(log, httpclient.New(service.ServiceNameStockbit))

	return &getTradeBook{
		base: base{
			logger:       log,
			emittenStore: emittenStore,
			symbols:      parseSymbols(symbols),
		},
		stockbit: stockbit,

		ctx:    ctx,
		cancel: cancel,
	}
}

type getTradeBook struct {
	base
	stockbit service.Stockbit

	ctx    context.Context
	cancel context.CancelFunc
}

func (u *getTradeBook) Run() error {
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

			tradeBook, err := stockbit.Retryable(func() (*service.TradeBook, error) {
				return u.stockbit.GetTradeBook(ctx, emitten)
			}).WithRetry()
			if err != nil {
				errs <- fmt.Errorf("failed to get trade book %s: %w", emitten, err)
				return
			}

			// TODO: analyze trade book & send to telegram if any anomaly

			u.logger.Info("successfully upserted broker summary", "symbol", emitten, "tradeBook", tradeBook)
		}(emittens[i])
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("failed to upsert broker summary: %w", errors.Join(<-errs))
	}

	return nil
}

func (u *getTradeBook) Stop() error {
	u.cancel()
	return nil
}
