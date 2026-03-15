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
		logger:       log,
		stockbit:     stockbit,
		emittenStore: emittenStore,
		ctx:          ctx,
		cancel:       cancel,
		symbols:      parseSymbols(symbols),
	}
}

type getTradeBook struct {
	logger   *slog.Logger
	stockbit service.Stockbit

	emittenStore service.EmittenStore

	ctx     context.Context
	cancel  context.CancelFunc
	symbols []string
}

func (u *getTradeBook) Run() (err error) {
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
