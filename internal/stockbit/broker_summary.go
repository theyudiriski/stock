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

var (
	investorTypeMap = map[service.InvestorType]string{
		service.InvestorTypeDomestic: "INVESTOR_TYPE_DOMESTIC",
		service.InvestorTypeForeign:  "INVESTOR_TYPE_FOREIGN",
	}

	marketBoardMap = map[service.MarketBoard]string{
		service.MarketBoardRegular:   "MARKET_BOARD_REGULER",
		service.MarketBoardNegosiasi: "MARKET_BOARD_NEGO",
		service.MarketBoardTunai:     "MARKET_BOARD_TUNAI",
	}
)

const (
	TransactionTypeGross = "TRANSACTION_TYPE_GROSS"
)

func (s *stockbit) GetBrokerSummary(
	ctx context.Context,
	symbol string,
	summaryDate string,
	investorType service.InvestorType,
	marketBoard service.MarketBoard,
) (
	*service.MarketDetector, error,
) {
	broksumURL, err := url.JoinPath(s.config.BaseURL, "marketdetectors", symbol)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(broksumURL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()
	query.Set("from", summaryDate)
	query.Set("to", summaryDate)
	query.Set("transaction_type", TransactionTypeGross)
	query.Set("market_board", marketBoardMap[marketBoard])
	query.Set("investor_type", investorTypeMap[investorType])
	query.Set("limit", "25")
	parsedURL.RawQuery = query.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	uniqueHash := fmt.Sprintf("%s-%s-%s-%s", symbol, summaryDate, investorType, marketBoard)
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

	var marketDetector service.MarketDetector
	err = json.Unmarshal(body, &marketDetector)
	if err != nil {
		return nil, err
	}

	// Inject InvestorType and Action for Buy brokers
	for i := range marketDetector.Data.BrokerSummary.BrokersBuy {
		marketDetector.Data.BrokerSummary.BrokersBuy[i].InvestorType = investorType
		marketDetector.Data.BrokerSummary.BrokersBuy[i].Action = service.ActionBuy
		marketDetector.Data.BrokerSummary.BrokersBuy[i].MarketBoard = marketBoard
	}

	// Inject InvestorType and Action for Sell brokers
	for i := range marketDetector.Data.BrokerSummary.BrokersSell {
		marketDetector.Data.BrokerSummary.BrokersSell[i].InvestorType = investorType
		marketDetector.Data.BrokerSummary.BrokersSell[i].Action = service.ActionSell
		marketDetector.Data.BrokerSummary.BrokersSell[i].MarketBoard = marketBoard
	}

	return &marketDetector, nil
}
