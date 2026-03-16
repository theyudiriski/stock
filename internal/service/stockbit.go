package service

import (
	"context"
)

type Stockbit interface {
	GetEmittenProfile(ctx context.Context, symbol string) (*EmittenProfile, error)
	GetBrokerSummary(ctx context.Context, symbol string, summaryDate string, investorType InvestorType, marketBoard MarketBoard) (*MarketDetector, error)

	GetPriceFeed(ctx context.Context, symbol string, fromDate string, toDate string) ([]PriceFeed, error)
	GetPriceIntraday(ctx context.Context, symbol string, date string) ([]PriceIntraday, error)

	GetRunningTrade(ctx context.Context, symbol, summaryDate string, tradeNumber *string) (*RunningTrade, error)
	GetSubsectors(ctx context.Context, sectorID string) ([]Sector, error)
	GetEmittenInfo(ctx context.Context, symbol string) (*EmittenInfo, error)

	GetDividends(ctx context.Context, symbol string) (*[]Dividend, error)
	GetRUPS(ctx context.Context, symbol string) (*[]RUPS, error)
	GetPublicExpose(ctx context.Context, symbol string) (*[]PublicExpose, error)
	GetRightIssue(ctx context.Context, symbol string) (*[]RightIssue, error)
	GetStockSplit(ctx context.Context, symbol string) (*[]Split, error)
	GetReverseSplit(ctx context.Context, symbol string) (*[]Split, error)
	// GetCorporateActions(ctx context.Context, symbol string) (*Corpaction, error)

	GetSubsidiaryCompanies(ctx context.Context, symbol string) (*SubsidiaryData, error)

	GetShareholders(ctx context.Context, symbol string, timeframe *Timeframe) (*ShareholderChartData, error)
	GetTradeBook(ctx context.Context, symbol string) (*TradeBook, error)
}
