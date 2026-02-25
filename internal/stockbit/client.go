package stockbit

import (
	"hash/fnv"
	"net/http"
	"stock/config"
	"stock/internal/service"
)

func NewStockbit(
	client *http.Client,
) service.Stockbit {
	cfg := config.LoadStockbit()
	return &stockbit{
		client:     client,
		config:     cfg,
		tokenCount: uint64(len(cfg.Tokens)),
	}
}

type stockbit struct {
	client     *http.Client
	config     config.Stockbit
	tokenCount uint64
}

// getTokenByID returns a token based on deterministic hashing of the identifier.
func (s *stockbit) getToken(hash string) string {
	if s.tokenCount == 1 {
		return s.config.Tokens[0]
	}

	h := fnv.New64a()
	h.Write([]byte(hash))
	idx := h.Sum64() % s.tokenCount

	return s.config.Tokens[idx]
}
