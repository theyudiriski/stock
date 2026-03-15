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
	"time"
)

func NewUpsertBrokerSummary(fromDate, toDate, symbols string) *upsertBrokerSummary {
	ctx, cancel := context.WithCancel(context.Background())

	logger.Init()
	log := logger.Default

	db, err := postgres.NewClient(config.LoadDatabase(), false)
	if err != nil {
		panic(err)
	}
	emittenStore := postgres.NewEmittenStore(db)
	brokerSummaryStore := postgres.NewBrokerSummaryStore(db)

	stockbit := stockbit.NewStockbit(log, httpclient.New(service.ServiceNameStockbit))

	return &upsertBrokerSummary{
		logger:             log,
		stockbit:           stockbit,
		emittenStore:       emittenStore,
		brokerSummaryStore: brokerSummaryStore,
		ctx:                ctx,
		cancel:             cancel,

		fromDate: fromDate,
		toDate:   toDate,
		symbols:  parseSymbols(symbols),
	}
}

type upsertBrokerSummary struct {
	logger             *slog.Logger
	stockbit           service.Stockbit
	emittenStore       service.EmittenStore
	brokerSummaryStore service.BrokerSummaryStore

	ctx    context.Context
	cancel context.CancelFunc

	fromDate string
	toDate   string
	symbols  []string
}

type brokerSummary struct {
	foreign  *service.MarketDetector
	domestic *service.MarketDetector
}

func (u *upsertBrokerSummary) Run() (err error) {
	start := time.Now()
	ctx := u.ctx

	emittens := u.symbols
	if len(emittens) < 1 {
		emittens, err = u.emittenStore.GetEmittens(ctx)
		if err != nil {
			u.logger.Error("failed to get emittens", "error", err)
			return err
		}
	}

	from, err := time.Parse(time.DateOnly, u.fromDate)
	if err != nil {
		return err
	}
	to, err := time.Parse(time.DateOnly, u.toDate)
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

			for date := from; !date.After(to); date = date.AddDate(0, 0, 1) {
				if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
					summaryDate := date.Format(time.DateOnly)

					result, err := stockbit.Retryable(func() (*brokerSummary, error) {
						foreignSummary, err := u.stockbit.GetBrokerSummary(ctx, emitten, summaryDate, service.InvestorTypeForeign, service.MarketBoardRegular)
						if err != nil {
							return nil, err
						}

						domesticSummary, err := u.stockbit.GetBrokerSummary(ctx, emitten, summaryDate, service.InvestorTypeDomestic, service.MarketBoardRegular)
						if err != nil {
							return nil, err
						}

						return &brokerSummary{foreign: foreignSummary, domestic: domesticSummary}, nil
					}).WithRetry()
					if err != nil {
						u.logger.Warn("failed to get broker summary", "symbol", emitten, "date", summaryDate, "error", err)
						errs <- err
						return
					}

					summary := &service.MarketDetector{
						Data: service.MarketDetectorData{
							BrokerSummary: service.BrokerSummary{
								BrokersBuy:  append(result.foreign.Data.BrokerSummary.BrokersBuy, result.domestic.Data.BrokerSummary.BrokersBuy...),
								BrokersSell: append(result.foreign.Data.BrokerSummary.BrokersSell, result.domestic.Data.BrokerSummary.BrokersSell...),
							},
						},
					}

					if len(summary.Data.BrokerSummary.BrokersBuy) > 0 || len(summary.Data.BrokerSummary.BrokersSell) > 0 {
						err = u.brokerSummaryStore.UpsertBrokerSummary(ctx, emitten, summaryDate, summary)
						if err != nil {
							u.logger.Warn("failed to upsert broker summary", "symbol", emitten, "date", summaryDate, "error", err)
							errs <- err
							return
						}
					}
				}
			}

			u.logger.Info("successfully upserted broker summary", "symbol", emitten)
		}(emittens[i])
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("failed to upsert broker summary: %w", errors.Join(<-errs))
	}

	u.logger.Info("successfully upserted broker summary", "duration", time.Since(start), "fromDate", u.fromDate, "toDate", u.toDate)
	return nil
}

func (u *upsertBrokerSummary) Stop() error {
	u.cancel()
	return nil
}
