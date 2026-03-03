package service

import "context"

type ShareholderStore interface {
	GetShareholderChartHistory(ctx context.Context, symbol string) (lastUpdatedDate *string, err error)
	UpsertShareholderChart(ctx context.Context, symbol string, data *ShareholderChartData) error
}

type InvestorCategory string

const (
	InvestorCategoryIndividual           InvestorCategory = "ID"
	InvestorCategoryCorporate            InvestorCategory = "CP"
	InvestorCategoryMutualFund           InvestorCategory = "MF"
	InvestorCategoryFinancialInstitution InvestorCategory = "IB"
	InvestorCategoryInsurance            InvestorCategory = "IS"
	InvestorCategorySecuritiesCompany    InvestorCategory = "SC"
	InvestorCategoryPensionFund          InvestorCategory = "PF"
	InvestorCategoryFoundation           InvestorCategory = "FD"
	InvestorCategoryOthers               InvestorCategory = "OT"
)

type (
	ShareholderType string
	Timeframe       int
)

const (
	ShareholderTypeForeign ShareholderType = "foreign"
	ShareholderTypeLocal   ShareholderType = "local"

	TimeframeFiveMonths Timeframe = 5
	TimeframeThreeYears Timeframe = 36
)

type ShareholderChartResponse struct {
	Message string               `json:"message"`
	Data    ShareholderChartData `json:"data"`
}

type ShareholderChartData struct {
	LastUpdate  string             `json:"last_update"`
	Shareholder []ShareholderChart `json:"legend"`
}

type ShareholderChart struct {
	ItemName     string                 `json:"item_name"`
	InvestorType InvestorType           `json:"-"`
	ChartData    []ShareholderChartItem `json:"chart_data"`
}

type ShareholderChartItem struct {
	Value    float64 `json:"value"`
	UnixDate string  `json:"unix_date"`
}
