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

func (s *stockbit) GetDividends(ctx context.Context, symbol string) (*[]service.Dividend, error) {
	dividendURL, err := url.JoinPath(s.config.BaseURL, "corpaction", "dividend")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(dividendURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
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

	var result service.CorpactionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.Dividends, nil
}

func (s *stockbit) GetRUPS(ctx context.Context, symbol string) (*[]service.RUPS, error) {
	rupsURL, err := url.JoinPath(s.config.BaseURL, "corpaction", "rups")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(rupsURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
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

	var result service.CorpactionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.RUPS, nil
}

func (s *stockbit) GetPublicExpose(ctx context.Context, symbol string) (*[]service.PublicExpose, error) {
	pubexURL, err := url.JoinPath(s.config.BaseURL, "corpaction", "pubex")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(pubexURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
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

	var result service.CorpactionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.PublicExpose, nil
}

func (s *stockbit) GetRightIssue(ctx context.Context, symbol string) (*[]service.RightIssue, error) {
	rightissueURL, err := url.JoinPath(s.config.BaseURL, "corpaction", "rightissue")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(rightissueURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
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

	var result service.CorpactionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.RightIssue, nil
}

func (s *stockbit) GetStockSplit(ctx context.Context, symbol string) (*[]service.Split, error) {
	stockSplitURL, err := url.JoinPath(s.config.BaseURL, "corpaction", "stocksplit")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(stockSplitURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
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

	var result service.CorpactionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Data.StockSplit, nil
}

func (s *stockbit) GetReverseSplit(ctx context.Context, symbol string) (*[]service.Split, error) {
	reverseSplitURL, err := url.JoinPath(s.config.BaseURL, "corpaction", "reversesplit")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(reverseSplitURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
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

	var result service.CorpactionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Data.ReverseSplit, nil
}

func (s *stockbit) GetCorporateActions(ctx context.Context, symbol string) (*service.Corpaction, error) {
	uri, err := url.JoinPath(s.config.BaseURL, "corpaction", symbol)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("limit", "5")
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
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

	var result service.CorpactionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}
