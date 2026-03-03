package service

import "context"

type CorpactionStore interface {
	UpsertCorpactions(ctx context.Context, symbol string, corpaction *Corpaction) error
}

type CorpactionResponse struct {
	Data Corpaction `json:"data"`
}

type Corpaction struct {
	Dividends    *[]Dividend     `json:"dividend"`
	RUPS         *[]RUPS         `json:"rups"`
	PublicExpose *[]PublicExpose `json:"pubex"`
	RightIssue   *[]RightIssue   `json:"rightissue"`
	StockSplit   *[]Split        `json:"stocksplit"`
	ReverseSplit *[]Split        `json:"stock_reverse"`
}

type Dividend struct {
	CumDate       string `json:"dividend_cumdate"`
	ExDate        string `json:"dividend_exdate"`
	PaymentDate   string `json:"dividend_paydate"`
	RecordingDate string `json:"dividend_recdate"`
	Amount        string `json:"dividend_value"`
	Currency      string `json:"dividend_currency"`
}

type RUPS struct {
	Date         string `json:"rups_date"` // 2025-06-25
	Time         string `json:"rups_time"` // 13:00
	Venue        string `json:"rups_venue"`
	EligibleDate string `json:"rups_eligible_date"`
}

type PublicExpose struct {
	Date  string `json:"puexp_date"` // 2025-09-29
	Time  string `json:"puexp_time"` // 14:00:00
	Venue string `json:"puexp_venue"`
}

type RightIssue struct {
	CumDate       string `json:"rightissue_cumdate"`
	ExDate        string `json:"rightissue_exdate"`
	RecordingDate string `json:"rightissue_recdate"`
	StartDate     string `json:"rightissue_trading_end"`
	EndDate       string `json:"rightissue_trading_start"`
	Price         int    `json:"rightissue_price"`
	Old           string `json:"rightissue_old"`
	New           string `json:"rightissue_new"`
}

type Split struct {
	CumDate       string `json:"stocksplit_cumdate"`
	ExDate        string `json:"stocksplit_exdate"`
	RecordingDate string `json:"stocksplit_recdate"`
	Old           string `json:"stocksplit_old"`
	New           string `json:"stocksplit_new"`
}
