package cronjob

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"stock/config"
	"stock/internal/postgres"
	"stock/internal/service"
	"stock/internal/stockbit"
	"strings"
	"sync"
)

func NewUpsertPriceFeed(fromDate, toDate, symbols string) *upsertPriceFeed {
	ctx, cancel := context.WithCancel(context.Background())

	stockbit := stockbit.NewStockbit(http.DefaultClient)

	db, err := postgres.NewClient(config.LoadDatabase(), false)
	if err != nil {
		panic(err)
	}
	emittenStore := postgres.NewEmittenStore(db)
	priceFeedStore := postgres.NewPriceFeedStore(db)

	return &upsertPriceFeed{
		stockbit:       stockbit,
		emittenStore:   emittenStore,
		priceFeedStore: priceFeedStore,
		ctx:            ctx,
		cancel:         cancel,

		fromDate: fromDate,
		toDate:   toDate,
		symbols:  strings.Split(symbols, ","),
	}
}

type upsertPriceFeed struct {
	stockbit       service.Stockbit
	emittenStore   service.EmittenStore
	priceFeedStore service.PriceFeedStore
	ctx            context.Context
	cancel         context.CancelFunc

	fromDate string
	toDate   string
	symbols  []string
}

func (u *upsertPriceFeed) Run() (err error) {
	ctx := u.ctx

	emittens := u.symbols
	if len(emittens) < 1 {
		emittens, err = u.emittenStore.GetEmittens(ctx)
		if err != nil {
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
		}(emittens[i])
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("failed to upsert price feed: %w", errors.Join(<-errs))
	}

	return nil
}

func (u *upsertPriceFeed) Stop() error {
	return nil
}
