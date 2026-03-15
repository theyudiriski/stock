package stockbit

import (
	"hash/fnv"
	"log/slog"
	"net/http"
	"sort"
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
		tokens:            getTokens(cfg.Tokens),
		tokenCount:        uint64(len(cfg.Tokens)),
		webviewTokenCount: uint64(len(cfg.WebviewTokens)),
	}
}

func getTokens(m map[string]string) []token {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tokens := make([]token, 0, len(m))
	for _, k := range keys {
		tokens = append(tokens, token{token: k, username: m[k]})
	}
	return tokens
}

type token struct {
	token    string
	username string
}

type stockbit struct {
	logger            *slog.Logger
	client            *http.Client
	config            config.Stockbit
	tokens            []token
	tokenCount        uint64
	webviewTokenCount uint64
}

// getToken returns a token and its associated username based on deterministic hashing of the identifier.
func (s *stockbit) getToken(hash string) (token, username string) {
	var idx uint64
	if s.tokenCount > 1 {
		h := fnv.New64a()
		h.Write([]byte(hash))
		idx = h.Sum64() % s.tokenCount
	}
	entry := s.tokens[idx]
	return entry.username, entry.token
}

// getWebviewToken returns a webview token based on deterministic hashing of the identifier.
func (s *stockbit) getWebviewToken(hash string) string {
	var idx uint64
	if s.webviewTokenCount > 1 {
		h := fnv.New64a()
		h.Write([]byte(hash))
		idx = h.Sum64() % s.webviewTokenCount
	}
	return s.config.WebviewTokens[idx]
}
