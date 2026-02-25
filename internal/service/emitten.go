package service

import (
	"context"
)

type EmittenStore interface {
	GetEmittens(ctx context.Context) ([]string, error)
	UpsertEmittenProfile(ctx context.Context, symbol string, profile *EmittenProfile) error
	UpsertEmittenProfileInfo(ctx context.Context, symbol string, info *EmittenInfo, profile *EmittenProfile) error

	GetEmittensUnderwriters(ctx context.Context) ([]EmittenUnderwriters, error)
	UpdateEmittenUnderwriterCode(ctx context.Context, symbol string, brokers []Broker) error
}

type EmittenProfileResponse struct {
	Data    EmittenProfile `json:"data"`
	Message string         `json:"message"`
}

type EmittenProfile struct {
	Description                     string              `json:"background"`
	History                         History             `json:"history"`
	Subsidiary                      []Subsidiary        `json:"subsidiary"`
	Shareholder                     []Shareholder       `json:"shareholder"`
	ShareholderDirectorCommissioner []Shareholder       `json:"shareholder_director_commissioner"`
	ShareholderNumbers              []ShareholderNumber `json:"shareholder_numbers"`
}

type History struct {
	Underwriters []string `json:"underwriters"`
	FreeFloat    string   `json:"free_float"`
}

type Shareholder struct {
	Percentage string   `json:"percentage"`
	Name       string   `json:"name"`
	Value      string   `json:"value"`
	Badges     []string `json:"badges"`
}

type Subsidiary struct {
	Company    string `json:"company"`
	Percentage string `json:"percentage"`
}

type ShareholderNumber struct {
	ShareholderDate string `json:"shareholder_date"`
	TotalShare      string `json:"total_share"`
	Change          int    `json:"change"`
	ChangeFormatted string `json:"change_formatted"`
}

type EmittenInfoResponse struct {
	Data    EmittenInfo `json:"data"`
	Message string      `json:"message"`
}

type EmittenInfo struct {
	Name     string    `json:"name"`
	Catalogs []Catalog `json:"catalogs"`
}

type Catalog struct {
	ID          string `json:"id"`
	CompanyType string `json:"company_type"`
}

type EmittenUnderwriters struct {
	Symbol       string   `json:"symbol"`
	Underwriters []string `json:"underwriters"`
}
