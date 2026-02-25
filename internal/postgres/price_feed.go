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
	priceFeed *service.PriceFeed,
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

	for i := range priceFeed.Data.Result {
		_, err = stmt.ExecContext(ctx,
			symbol,
			priceFeed.Data.Result[i].Date,
			priceFeed.Data.Result[i].Open,
			priceFeed.Data.Result[i].Close,
			priceFeed.Data.Result[i].High,
			priceFeed.Data.Result[i].Low,
			priceFeed.Data.Result[i].Average,
			priceFeed.Data.Result[i].Value,
			priceFeed.Data.Result[i].Volume,
			priceFeed.Data.Result[i].Frequency,
			priceFeed.Data.Result[i].NetForeignBuy,
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
