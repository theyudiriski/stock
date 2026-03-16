package postgres

import (
	"context"
	"fmt"
	"stock/internal/service"
	"strconv"
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
			transaction_value,
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

func (s *priceFeedStore) UpsertPriceIntraday(
	ctx context.Context,
	symbol string,
	priceIntradays []service.PriceIntraday,
) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO price_intradays (
			symbol,
			unix_timestamp,
			open,
			close,
			high,
			low,
			transaction_value,
			volume
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (symbol, unix_timestamp)
		DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range priceIntradays {
		unixTimestamp, err := strconv.ParseInt(priceIntradays[i].UnixTimestamp, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse unix timestamp: %w", err)
		}

		volume, err := strconv.ParseInt(priceIntradays[i].Volume, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse volume: %w", err)
		}

		_, err = stmt.ExecContext(ctx,
			symbol,
			unixTimestamp,
			priceIntradays[i].Open,
			priceIntradays[i].Close,
			priceIntradays[i].High,
			priceIntradays[i].Low,
			priceIntradays[i].Value,
			volume,
		)
		if err != nil {
			return fmt.Errorf("failed to insert price intraday: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
