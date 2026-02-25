package service

import "context"

type SectorStore interface {
	UpsertSubsectors(ctx context.Context, sectors []Sector) error
}

type SectorResponse struct {
	Data    []Sector `json:"data"`
	Message string   `json:"message"`
}

type Sector struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent"`
}
