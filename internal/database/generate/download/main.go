package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gabe565/gones/internal/log"
)

func main() {
	log.Init(os.Stderr)

	action, err := NewDownloader("Nintendo - Nintendo Entertainment System")
	if err != nil {
		slog.Error("Failed to create downloader", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := action.Run(ctx); err != nil {
		slog.Error("Failed to run downloader", "error", err)
		os.Exit(1) //nolint:gocritic
	}
}
