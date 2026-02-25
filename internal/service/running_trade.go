package service

import (
	"context"
)

type RunningTradeStore interface {
	UpsertRunningTrade(ctx context.Context, symbol string, summaryDate string, summary *RunningTrade) error
}

type RunningTrade struct {
	Data    RunningTradeData `json:"data"`
	Message string           `json:"message"`
}

type RunningTradeData struct {
	IsOpenMarket bool               `json:"is_open_market"`
	Date         string             `json:"date"`
	Items        []RunningTradeItem `json:"running_trade"`
}

type RunningTradeItem struct {
	Time        string      `json:"time"`
	Action      string      `json:"action"`
	Price       string      `json:"price"`
	Lot         string      `json:"lot"`
	Buyer       string      `json:"buyer"`
	Seller      string      `json:"seller"`
	TradeNumber string      `json:"trade_number"`
	MarketBoard MarketBoard `json:"market_board"`
}
