package service

import "context"

type BrokerStore interface {
	GetBrokerByName(ctx context.Context, name string) (*Broker, error)
}

type Broker struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
