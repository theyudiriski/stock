package logger

import (
	"log/slog"
	"os"
	"strings"

	"stock/config"
)

var Default *slog.Logger

func Init() {
	cfg := config.LoadLogger()

	level := parseLevel(cfg.Level)
	format := strings.ToLower(cfg.Format)

	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if format == "json" {
		h = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		h = slog.NewTextHandler(os.Stderr, opts)
	}

	Default = slog.New(h)
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
