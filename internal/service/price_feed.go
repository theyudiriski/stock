package service

import "context"

type PriceFeedStore interface {
	UpsertPriceFeed(ctx context.Context, symbol string, priceFeed *PriceFeed) error
}

type PriceFeed struct {
	Data    PriceFeedData `json:"data"`
	Message string        `json:"message"`
}

type PriceFeedData struct {
	Result   []PriceFeedResult `json:"result"`
	Paginate Paginate          `json:"paginate"`
}

type Paginate struct {
	NextPage string `json:"next_page"`
}

type PriceFeedResult struct {
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
