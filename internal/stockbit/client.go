package stockbit

import (
	"hash/fnv"
	"log/slog"
	"net/http"
	"stock/config"
	"stock/internal/service"
)

func NewStockbit(
	logger *slog.Logger,
	client *http.Client,
) service.Stockbit {
	cfg := config.LoadStockbit()
	return &stockbit{
		logger:            logger,
		client:            client,
		config:            cfg,
		tokenCount:        uint64(len(cfg.Tokens)),
		webviewTokenCount: uint64(len(cfg.WebviewTokens)),
	}
}

type stockbit struct {
	logger            *slog.Logger
	client            *http.Client
	config            config.Stockbit
	tokenCount        uint64
	webviewTokenCount uint64
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

// getWebviewToken returns a webview token based on deterministic hashing of the identifier.
func (s *stockbit) getWebviewToken(hash string) string {
	if s.tokenCount == 1 {
		return s.config.WebviewTokens[0]
	}

	h := fnv.New64a()
	h.Write([]byte(hash))
	idx := h.Sum64() % s.webviewTokenCount

	return s.config.WebviewTokens[idx]
}
