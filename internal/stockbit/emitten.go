package stockbit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"stock/internal/service"
	"time"
)

func (s *stockbit) GetEmittenProfile(ctx context.Context, symbol string) (*service.EmittenProfile, error) {
	emittenProfileURL, err := url.JoinPath(s.config.BaseURL, "emitten", symbol, "profile")
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, emittenProfileURL, nil)
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
		return nil, fmt.Errorf("failed to get emitten profile: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var emittenProfile service.EmittenProfileResponse
	err = json.Unmarshal(body, &emittenProfile)
	if err != nil {
		return nil, err
	}

	// loop through emitten beneficiaries to preserve uniqueness
	beneficiaries := make([]service.Beneficiary, 0)
	for i := range emittenProfile.Data.Beneficiary {
		if !slices.Contains(beneficiaries, emittenProfile.Data.Beneficiary[i]) {
			beneficiaries = append(beneficiaries, emittenProfile.Data.Beneficiary[i])
		}
	}
	emittenProfile.Data.Beneficiary = beneficiaries

	return &emittenProfile.Data, nil
}

func (s *stockbit) GetEmittenInfo(ctx context.Context, symbol string) (*service.EmittenInfo, error) {
	emittenInfoURL, err := url.JoinPath(s.config.BaseURL, "emitten", symbol, "info")
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, emittenInfoURL, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", "Bearer "+s.getToken(symbol))

	response, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get emitten info: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var emittenInfo service.EmittenInfoResponse
	err = json.Unmarshal(body, &emittenInfo)
	if err != nil {
		return nil, err
	}

	return &emittenInfo.Data, nil
}
