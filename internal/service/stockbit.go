package service

import (
	"context"
)

type Stockbit interface {
	GetEmittenProfile(ctx context.Context, symbol string) (*EmittenProfile, error)
	GetBrokerSummary(ctx context.Context, symbol string, summaryDate string, investorType InvestorType, marketBoard MarketBoard) (*MarketDetector, error)
	GetPriceFeed(ctx context.Context, symbol string, fromDate string, toDate string) (*PriceFeed, error)
	GetRunningTrade(ctx context.Context, symbol, summaryDate string, tradeNumber *string) (*RunningTrade, error)
	GetSubsectors(ctx context.Context, sectorID string) ([]Sector, error)
	GetEmittenInfo(ctx context.Context, symbol string) (*EmittenInfo, error)
}
