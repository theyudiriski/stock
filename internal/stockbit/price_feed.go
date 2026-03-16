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

func (s *stockbit) GetPriceFeed(
	ctx context.Context,
	symbol string,
	fromDate string,
	toDate string,
) ([]service.PriceFeed, error) {
	pricefeedURL, err := url.JoinPath(s.config.BaseURL, "/company-price-feed/historical/summary", symbol)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(pricefeedURL)
	if err != nil {
		return nil, err
	}

	var (
		results []service.PriceFeed
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
		token, username := s.getToken(uniqueHash)
		r.Header.Set("Authorization", "Bearer "+token)

		response, err := s.client.Do(r)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return nil, s.handleError(response, username)
		}

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return nil, err
		}

		var result service.PriceFeedResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}

		if len(result.Data.Result) < 1 {
			break
		}

		results = append(results, result.Data.Result...)
		page = result.Data.Paginate.NextPage
	}

	return results, nil
}

func (s *stockbit) GetPriceIntraday(
	ctx context.Context,
	symbol string,
	date string,
) ([]service.PriceIntraday, error) {
	uri, err := url.JoinPath(s.config.BaseURL, "chartbit", symbol, "price", "intraday")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return nil, err
	}

	from := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 0, loc)
	to := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 8, 58, 0, 0, loc)

	query := parsedURL.Query()
	query.Set("from", fmt.Sprintf("%d", from.Unix()))
	query.Set("to", fmt.Sprintf("%d", to.Unix()))
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%s", symbol, date)
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
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	var result service.PriceIntradayResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Chartbit, nil
}
