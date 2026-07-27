package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Dmitro44/relper/internal/config"
	"github.com/Dmitro44/relper/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	_, err = database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
}
