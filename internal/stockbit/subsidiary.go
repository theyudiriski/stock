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

func (s *stockbit) GetSubsidiaryCompanies(ctx context.Context, symbol string) (*service.SubsidiaryData, error) {
	uri, err := url.JoinPath(s.config.BaseURL, "emitten-metadata", "subsidiary", symbol)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%d", symbol, time.Now().UnixMilli())
	r.Header.Set("Authorization", "Bearer "+s.getToken(uniqueHash))

	response, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get subsidiary companies: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result service.SubsidiaryResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Data.Subsidiaries) < 1 && result.Data.LastUpdatedPeriod == "" {
		return nil, nil
	}

	return &result.Data, nil
}
