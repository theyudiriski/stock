package service

import (
	"context"
)

type SubsidiaryStore interface {
	GetSubsidiaryCompanyHistory(ctx context.Context, symbol string) (lastUpdatedPeriod *string, err error)
	UpsertSubsidiaryCompanies(ctx context.Context, symbol string, data *SubsidiaryData) error
}

type SubsidiaryResponse struct {
	Data    SubsidiaryData `json:"data"`
	Message string         `json:"message"`
}

type SubsidiaryData struct {
	Subsidiaries      []SubsidiaryCompany `json:"subsidiaries"`
	LastUpdatedPeriod string              `json:"last_updated_period"`
}

type SubsidiaryCompany struct {
	CompanyName  string `json:"company_name"`
	Percentage   string `json:"percentage"`
	BusinessType string `json:"business_type"`
}
