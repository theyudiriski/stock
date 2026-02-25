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

func (s *stockbit) GetPriceFeed(
	ctx context.Context,
	symbol string,
	fromDate string,
	toDate string,
) (*service.PriceFeed, error) {
	pricefeedURL, err := url.JoinPath(s.config.BaseURL, "/company-price-feed/historical/summary", symbol)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(pricefeedURL)
	if err != nil {
		return nil, err
	}

	var (
		results []service.PriceFeedResult
		page    = "1"
	)

	for {
		query := parsedURL.Query()
		query.Set("period", "HS_PERIOD_DAILY")
		query.Set("start_date", fromDate)
		query.Set("end_date", toDate)
		query.Set("limit", "50")
		query.Set("page", page)
		parsedURL.RawQuery = query.Encode()

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
		if err != nil {
			return nil, err
		}

		uniqueHash := fmt.Sprintf("%s-%s-%s-%s", symbol, fromDate, toDate, page)

		r.Header.Set("Authorization", "Bearer "+s.getToken(uniqueHash))

		response, err := s.client.Do(r)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return nil, fmt.Errorf("failed to get price feed: %d", response.StatusCode)
		}

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, err
		}

		var priceFeed service.PriceFeed
		err = json.Unmarshal(body, &priceFeed)
		if err != nil {
			return nil, err
		}

		if len(priceFeed.Data.Result) < 1 {
			break
		}

		results = append(results, priceFeed.Data.Result...)
		page = priceFeed.Data.Paginate.NextPage
	}

	return &service.PriceFeed{
		Data: service.PriceFeedData{
			Result: results,
		},
	}, nil
}
