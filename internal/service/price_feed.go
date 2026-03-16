package service

import "context"

type PriceFeedStore interface {
	UpsertPriceFeed(ctx context.Context, symbol string, priceFeed []PriceFeed) error
	UpsertPriceIntraday(ctx context.Context, symbol string, priceIntraday []PriceIntraday) error
}

type PriceFeedResponse struct {
	Data    PriceFeedData `json:"data"`
	Message string        `json:"message"`
}

type PriceFeedData struct {
	Result   []PriceFeed `json:"result"`
	Paginate Paginate    `json:"paginate"`
}

type Paginate struct {
	NextPage string `json:"next_page"`
}

type PriceFeed struct {
	Date          string `json:"date"`
	Open          int    `json:"open"`
	Close         int    `json:"close"`
	High          int    `json:"high"`
	Low           int    `json:"low"`
	Average       int    `json:"average"`
	Value         int64  `json:"value"`
	Volume        int64  `json:"volume"`
	Frequency     int64  `json:"frequency"`
	NetForeignBuy int64  `json:"net_foreign"`
}

type PriceIntradayResponse struct {
	Data    PriceIntradayData `json:"data"`
	Message string            `json:"message"`
}

type PriceIntradayData struct {
	Chartbit []PriceIntraday `json:"chartbit"`
}

type PriceIntraday struct {
	Open          int    `json:"open"`
	Close         int    `json:"close"`
	High          int    `json:"high"`
	Low           int    `json:"low"`
	UnixTimestamp string `json:"unix_timestamp"`
	Value         int64  `json:"value"`
	Volume        string `json:"volume"`
}
