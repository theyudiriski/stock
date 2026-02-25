package service

import "context"

type (
	InvestorType string
	MarketBoard  string
	Action       string
)

const (
	InvestorTypeDomestic InvestorType = "D"
	InvestorTypeForeign  InvestorType = "F"

	MarketBoardRegular   MarketBoard = "RG"
	MarketBoardNegosiasi MarketBoard = "NG"
	MarketBoardTunai     MarketBoard = "TN"

	ActionBuy  Action = "B"
	ActionSell Action = "S"
)

type BrokerSummaryStore interface {
	UpsertBrokerSummary(ctx context.Context, symbol, summaryDate string, summary *MarketDetector) error
}

type MarketDetector struct {
	Message string             `json:"message"`
	Data    MarketDetectorData `json:"data"`
}

type MarketDetectorData struct {
	BrokerSummary BrokerSummary `json:"broker_summary"`
	From          string        `json:"from"`
	To            string        `json:"to"`
}

type BrokerSummary struct {
	BrokersBuy  []BrokerBuy  `json:"brokers_buy"`
	BrokersSell []BrokerSell `json:"brokers_sell"`
}

type BrokerBuy struct {
	Blot             string       `json:"blot"`
	Blotv            string       `json:"blotv"`
	Bval             string       `json:"bval"`
	Bvalv            string       `json:"bvalv"`
	NetbsBrokerCode  string       `json:"netbs_broker_code"`
	NetbsBuyAvgPrice string       `json:"netbs_buy_avg_price"`
	NetbsDate        string       `json:"netbs_date"`
	NetbsStockCode   string       `json:"netbs_stock_code"`
	Type             string       `json:"type"`
	InvestorType     InvestorType `json:"-"` // D: Domestic, F: Foreign
	Action           Action       `json:"-"` // B: Buy, S: Sell
	MarketBoard      MarketBoard  `json:"-"` // RG: Regular, NG: Negosiasi, TN: Tunai
}

type BrokerSell struct {
	NetbsBrokerCode   string       `json:"netbs_broker_code"`
	NetbsDate         string       `json:"netbs_date"`
	NetbsSellAvgPrice string       `json:"netbs_sell_avg_price"`
	NetbsStockCode    string       `json:"netbs_stock_code"`
	Slot              string       `json:"slot"`
	Slotv             string       `json:"slotv"`
	Sval              string       `json:"sval"`
	Svalv             string       `json:"svalv"`
	Type              string       `json:"type"`
	InvestorType      InvestorType `json:"-"` // D: Domestic, F: Foreign
	Action            Action       `json:"-"` // B: Buy, S: Sell
	MarketBoard       MarketBoard  `json:"-"` // RG: Regular, NG: Negosiasi, TN: Tunai
}
