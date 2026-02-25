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

func (s *stockbit) GetSubsectors(
	ctx context.Context,
	sectorID string,
) ([]service.Sector, error) {
	sectorURL, err := url.JoinPath(s.config.BaseURL, "emitten", "sectors", sectorID, "subsectors")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(sectorURL)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", "Bearer "+s.getToken(sectorID))

	response, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get subsectors: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var subsectorResponse service.SectorResponse
	err = json.Unmarshal(body, &subsectorResponse)
	if err != nil {
		return nil, err
	}

	return subsectorResponse.Data, nil
}
