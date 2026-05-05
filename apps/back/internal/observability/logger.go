// Package observability provides structured logging, Prometheus metrics,
// and health-check utilities for the seonology-journey-back service.
package observability

import (
	"log/slog"
	"os"
)

// NewLogger returns a configured slog.Logger based on LOG_FORMAT env var.
// "json" produces JSON output; anything else produces text output.
func NewLogger() *slog.Logger {
	var handler slog.Handler

	level := parseLevel(os.Getenv("LOG_LEVEL"))

	switch os.Getenv("LOG_FORMAT") {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler)
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "warn", "WARN":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
