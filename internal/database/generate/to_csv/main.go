//go:build embed_nes_xml

package main

import (
	"compress/gzip"
	"encoding/csv"
	"io"
	"os"

	"github.com/gabe565/gones/internal/database/nointro"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	const path = "internal/database/database.csv"

	datafile, err := nointro.Load(nointro.Nes)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load NoInto database")
	}

	log.Info().Str("path", path).Msg("Creating CSV file")
	f, err := os.Create(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create CSV file")
	}

	log.Info().Str("path", path+".gz").Msg("Creating gzipped CSV file")
	gzf, err := os.Create(path + ".gz")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create gzipped CSV file")
	}
	gz := gzip.NewWriter(gzf)

	c := csv.NewWriter(io.MultiWriter(f, gz))
	log.Info().Int("count", len(datafile.Games)).Msg("Writing games to CSV")
	for _, game := range datafile.Games {
		for _, rom := range game.Roms {
			if err := c.Write([]string{rom.MD5, game.Name}); err != nil {
				log.Fatal().Err(err).Msg("Failed to write CSV")
			}
		}
	}
	c.Flush()

	log.Info().Msg("Closing files")
	if err := c.Error(); err != nil {
		log.Fatal().Err(err).Msg("Failed to write CSV")
	}

	if err := f.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close CSV file")
	}

	if err := gz.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close gzipped CSV writer")
	}
	if err := gzf.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close gzipped CSV file")
	}
}
