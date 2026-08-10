package logger

import (
	"log/slog"
	"os"
)

func New(env string) *slog.Logger {
	var handler slog.Handler
	level := slog.LevelDebug

	if env == "production" {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}
