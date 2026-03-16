package postgres

import (
	"context"
	"fmt"
	"stock/internal/service"
)

func NewPriceFeedStore(
	db *Client,
) service.PriceFeedStore {
	return &priceFeedStore{
		db: db,
	}
}

type priceFeedStore struct {
	db *Client
}

func (s *priceFeedStore) UpsertPriceFeed(
	ctx context.Context,
	symbol string,
	priceFeeds []service.PriceFeed,
) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO price_feeds (
			symbol,
			date,
			open,
			close,
			high,
			low,
			average,
			value,
			volume,
			frequency,
			net_foreign
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (symbol, date)
		DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range priceFeeds {
		_, err = stmt.ExecContext(ctx,
			symbol,
			priceFeeds[i].Date,
			priceFeeds[i].Open,
			priceFeeds[i].Close,
			priceFeeds[i].High,
			priceFeeds[i].Low,
			priceFeeds[i].Average,
			priceFeeds[i].Value,
			priceFeeds[i].Volume,
			priceFeeds[i].Frequency,
			priceFeeds[i].NetForeignBuy,
		)
		if err != nil {
			return fmt.Errorf("failed to insert price feed: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
