package stockbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stock/internal/service"
	"time"
)

var (
	loc, _ = time.LoadLocation("Asia/Jakarta")
)

// GetTradeBook can only get latest trade book when market is open
func (s *stockbit) GetTradeBook(
	ctx context.Context,
	symbol string,
) (*service.TradeBook, error) {
	uri, err := url.JoinPath(s.config.BaseURL, "/order-trade/trade-book")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	query.Set("group_by", "GROUP_BY_TIME")
	query.Set("time_interval", "10m")

	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().UTC().UnixMilli()
	uniqueHash := fmt.Sprintf("%s-%d", symbol, currentTime)
	token, username := s.getToken(uniqueHash)
	r.Header.Set("Authorization", "Bearer "+token)

	response, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, s.handleError(response, username)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var tradeBook service.TradeBookResponse
	err = json.Unmarshal(body, &tradeBook)
	if err != nil {
		return nil, err
	}

	return &tradeBook.Data, nil
}
