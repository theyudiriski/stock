package cronjob

import (
	"context"
	"log/slog"
	"stock/internal/service"
	"strings"
)

func parseSymbols(symbols string) []string {
	if symbols != "" {
		return strings.Split(symbols, ",")
	}
	return nil
}

type base struct {
	emittenStore service.EmittenStore
	symbols      []string
	logger       *slog.Logger
}

func (b *base) getEmittens(ctx context.Context) (emittens []string, err error) {
	emittens = b.symbols
	if len(emittens) < 1 {
		emittens, err = b.emittenStore.GetEmittens(ctx)
		if err != nil {
			b.logger.Error("failed to get emittens", "error", err)
			return nil, err
		}
	}
	return emittens, nil
}
