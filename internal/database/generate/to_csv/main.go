//go:build embed_nes_xml

package main

import (
	"compress/gzip"
	"encoding/csv"
	"io"
	"log/slog"
	"os"

	"gabe565.com/gones/internal/database/nointro"
	"gabe565.com/gones/internal/log"
)

func main() {
	log.Init(os.Stderr)

	const path = "internal/database/database.csv"

	datafile, err := nointro.Load(nointro.Nes)
	if err != nil {
		slog.Error("Failed to load NoInto database", "error", err)
		os.Exit(1)
	}

	slog.Info("Creating CSV file", "path", path)
	f, err := os.Create(path)
	if err != nil {
		slog.Error("Failed to create CSV file", "error", err)
		os.Exit(1)
	}

	slog.Info("Creating gzipped CSV file", "path", path+".gz")
	gzf, err := os.Create(path + ".gz")
	if err != nil {
		slog.Error("Failed to create gzipped CSV file", "error", err)
		os.Exit(1)
	}
	gz := gzip.NewWriter(gzf)

	c := csv.NewWriter(io.MultiWriter(f, gz))
	slog.Info("Writing games to CSV", "count", len(datafile.Games))
	for _, game := range datafile.Games {
		for _, rom := range game.Roms {
			if err := c.Write([]string{rom.MD5, game.Name}); err != nil {
				slog.Error("Failed to write CSV", "error", err)
				os.Exit(1)
			}
		}
	}
	c.Flush()

	slog.Info("Closing files")
	if err := c.Error(); err != nil {
		slog.Error("Failed to write CSV", "error", err)
		os.Exit(1)
	}

	if err := f.Close(); err != nil {
		slog.Error("Failed to close CSV file", "error", err)
		os.Exit(1)
	}

	if err := gz.Close(); err != nil {
		slog.Error("Failed to close gzipped CSV writer", "error", err)
		os.Exit(1)
	}
	if err := gzf.Close(); err != nil {
		slog.Error("Failed to close gzipped CSV file", "error", err)
		os.Exit(1)
	}
}
