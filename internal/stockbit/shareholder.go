package stockbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"stock/internal/service"
	"strconv"
	"time"
)

var (
	shareholderTypeMap = map[service.ShareholderType]service.InvestorType{
		service.ShareholderTypeForeign: service.InvestorTypeForeign,
		service.ShareholderTypeLocal:   service.InvestorTypeDomestic,
	}
)

func (s *stockbit) GetShareholders(
	ctx context.Context,
	symbol string,
	tf *service.Timeframe,
) (*service.ShareholderChartData, error) {
	uri, err := url.JoinPath(s.config.BaseURL, "emitten-metadata", "shareholders", symbol, "chart")
	if err != nil {
		return nil, err
	}

	// initialize with default timeframe
	timeframe := service.TimeframeFiveMonths
	if tf != nil {
		timeframe = *tf
	}

	var result service.ShareholderChartData

	for i, typ := range []service.ShareholderType{service.ShareholderTypeForeign, service.ShareholderTypeLocal} {
		parsedURL, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}

		query := parsedURL.Query()
		query.Set("value_year", strconv.Itoa(int(timeframe)))
		query.Set("shareholder_type", string(typ))
		parsedURL.RawQuery = query.Encode()

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
		if err != nil {
			return nil, err
		}

		uniqueHash := fmt.Sprintf("%s-%s-%d", symbol, typ, time.Now().UnixMilli())
		r.Header.Set("Authorization", s.getWebviewToken(uniqueHash))

		response, err := s.client.Do(r)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get shareholders: %d", response.StatusCode)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var resp service.ShareholderChartResponse
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, err
		}

		// parse last update
		var lastUpdateDate string
		lastUpdate, err := time.Parse("2 Jan 06", resp.Data.LastUpdate)
		if err != nil {
			s.logger.Error("failed to parse last update", "lastUpdate", resp.Data.LastUpdate, "error", err)
			return nil, err
		}
		lastUpdateDate = lastUpdate.Format("2006-01-02")

		// check if last update is the same
		if i == 0 {
			result.LastUpdate = lastUpdateDate
		} else {
			if lastUpdateDate != result.LastUpdate {
				return nil, fmt.Errorf("last_update mismatch: foreign=%q local=%q", result.LastUpdate, resp.Data.LastUpdate)
			}
		}

		// set investor type for each shareholder
		investorType := shareholderTypeMap[typ]
		for j := range resp.Data.Shareholder {
			resp.Data.Shareholder[j].InvestorType = investorType
		}
		result.Shareholder = append(result.Shareholder, resp.Data.Shareholder...)
	}

	return &result, nil
}
