package postgres

import (
	"context"
	"fmt"
	"stock/internal/service"
	"stock/internal/stockbit"
	"strings"

	"github.com/shopspring/decimal"
)

func NewRunningTradeStore(
	db *Client,
) service.RunningTradeStore {
	return &runningTradeStore{
		db: db,
	}
}

type runningTradeStore struct {
	db *Client
}

var (
	runningTradeActionMap = map[string]service.Action{
		stockbit.RunningTradeActionBuy:  service.ActionBuy,
		stockbit.RunningTradeActionSell: service.ActionSell,
	}
)

func (s *runningTradeStore) extractBroker(brokerCode string) (string, string) {
	return brokerCode[0:2], brokerCode[4:5]
}

func (s *runningTradeStore) UpsertRunningTrade(
	ctx context.Context,
	symbol string,
	summaryDate string,
	summary *service.RunningTrade,
) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO running_trades (
			symbol,
			date,
			buyer,
			buyer_investor_type,
			seller,
			seller_investor_type,
			market_board,
			action,
			price,
			lot,
			trade_number
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (symbol, date, trade_number)
		DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range summary.Data.Items {
		buyerBroker, buyerInvestorType := s.extractBroker(summary.Data.Items[i].Buyer)
		sellerBroker, sellerInvetorType := s.extractBroker(summary.Data.Items[i].Seller)

		lotDecimal, err := decimal.NewFromString(summary.Data.Items[i].Lot)
		if err != nil {
			return fmt.Errorf("failed to parse lot: %w", err)
		}

		priceDecimal, err := decimal.NewFromString(strings.ReplaceAll(summary.Data.Items[i].Price, ",", ""))
		if err != nil {
			return fmt.Errorf("failed to parse price: %w", err)
		}

		tradeNumberDecimal, err := decimal.NewFromString(summary.Data.Items[i].TradeNumber)
		if err != nil {
			return fmt.Errorf("failed to parse trade number: %w", err)
		}

		_, err = stmt.ExecContext(ctx,
			symbol,
			summaryDate,
			buyerBroker,
			buyerInvestorType,
			sellerBroker,
			sellerInvetorType,
			summary.Data.Items[i].MarketBoard,
			runningTradeActionMap[summary.Data.Items[i].Action],
			priceDecimal.IntPart(),
			lotDecimal.IntPart(),
			tradeNumberDecimal.IntPart(),
		)

		if err != nil {
			return fmt.Errorf("failed to insert running trade: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
