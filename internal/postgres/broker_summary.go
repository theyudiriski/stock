package postgres

import (
	"context"
	"fmt"
	"stock/internal/service"

	"github.com/shopspring/decimal"
)

func NewBrokerSummaryStore(
	db *Client,
) service.BrokerSummaryStore {
	return &brokerSummaryStore{
		db: db,
	}
}

type brokerSummaryStore struct {
	db *Client
}

func (s *brokerSummaryStore) UpsertBrokerSummary(
	ctx context.Context,
	symbol string,
	summaryDate string,
	summary *service.MarketDetector,
) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := range summary.Data.BrokerSummary.BrokersBuy {
		// Convert string values to int64 for total_lot and total_value
		totalLotDecimal, err := decimal.NewFromString(summary.Data.BrokerSummary.BrokersBuy[i].Blot)
		if err != nil {
			return fmt.Errorf("failed to parse total_lot for buy: %w", err)
		}

		totalValueDecimal, err := decimal.NewFromString(summary.Data.BrokerSummary.BrokersBuy[i].Bval)
		if err != nil {
			return fmt.Errorf("failed to parse total_value for buy: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO broker_summaries (
				symbol,
				broker,
				action,
				investor_type,
				market_board,
				summary_date,
				total_lot,
				total_value,
				price_average
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (symbol, broker, action, investor_type, market_board, summary_date)
			DO NOTHING;
		`,
			symbol,
			summary.Data.BrokerSummary.BrokersBuy[i].NetbsBrokerCode,
			summary.Data.BrokerSummary.BrokersBuy[i].Action,
			summary.Data.BrokerSummary.BrokersBuy[i].InvestorType,
			summary.Data.BrokerSummary.BrokersBuy[i].MarketBoard,
			summaryDate,
			totalLotDecimal.IntPart(),
			totalValueDecimal.IntPart(),
			summary.Data.BrokerSummary.BrokersBuy[i].NetbsBuyAvgPrice,
		)
		if err != nil {
			return fmt.Errorf("failed to insert broker summary: %w", err)
		}
	}

	for i := range summary.Data.BrokerSummary.BrokersSell {
		// Convert string values to int64 for total_lot and total_value
		totalLotDecimal, err := decimal.NewFromString(summary.Data.BrokerSummary.BrokersSell[i].Slot)
		if err != nil {
			return fmt.Errorf("failed to parse total_lot for sell: %w", err)
		}

		totalValueDecimal, err := decimal.NewFromString(summary.Data.BrokerSummary.BrokersSell[i].Sval)
		if err != nil {
			return fmt.Errorf("failed to parse total_value for sell: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO broker_summaries (
				symbol,
				broker,
				action,
				investor_type,
				market_board,
				summary_date,
				total_lot,
				total_value,
				price_average
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (symbol, broker, action, investor_type, market_board, summary_date)
			DO NOTHING;
		`,
			symbol,
			summary.Data.BrokerSummary.BrokersSell[i].NetbsBrokerCode,
			summary.Data.BrokerSummary.BrokersSell[i].Action,
			summary.Data.BrokerSummary.BrokersSell[i].InvestorType,
			summary.Data.BrokerSummary.BrokersSell[i].MarketBoard,
			summary.Data.BrokerSummary.BrokersSell[i].NetbsDate,
			totalLotDecimal.IntPart(),
			totalValueDecimal.IntPart(),
			summary.Data.BrokerSummary.BrokersSell[i].NetbsSellAvgPrice,
		)
		if err != nil {
			return fmt.Errorf("failed to insert broker summary: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
