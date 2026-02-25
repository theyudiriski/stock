package stockbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stock/internal/service"
)

const (
	RunningTradeActionBuy  = "buy"
	RunningTradeActionSell = "sell"

	RunningTradeOrderByTime = "RUNNING_TRADE_ORDER_BY_TIME"
)

func (s *stockbit) GetRunningTrade(
	ctx context.Context,
	symbol, summaryDate string,
	tradeNumber *string,
) (*service.RunningTrade, error) {
	runningTradeURL, err := url.JoinPath(s.config.BaseURL, "/order-trade/running-trade")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(runningTradeURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbols[]", symbol)
	query.Set("date", summaryDate)
	query.Set("sort", "ASC")
	query.Set("limit", "100")
	query.Set("order_by", RunningTradeOrderByTime)
	if tradeNumber != nil {
		query.Set("trade_number", *tradeNumber)
	}

	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%s", symbol, summaryDate)
	if tradeNumber != nil {
		uniqueHash += fmt.Sprintf("-%s", *tradeNumber)
	}

	r.Header.Set("Authorization", "Bearer "+s.getToken(uniqueHash))

	response, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get broker summary: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var runningTrade service.RunningTrade
	err = json.Unmarshal(body, &runningTrade)
	if err != nil {
		return nil, err
	}

	return &runningTrade, nil
}
