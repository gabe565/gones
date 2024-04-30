//go:build embed_nes_xml

package main

import (
	"encoding/csv"
	"os"

	"github.com/gabe565/gones/internal/database/nointro"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	datafile, err := nointro.Load(nointro.Nes)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load NoInto database")
	}

	f, err := os.Create("internal/database/database.csv")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create CSV file")
	}

	c := csv.NewWriter(f)
	for _, game := range datafile.Games {
		for _, rom := range game.Roms {
			if err := c.Write([]string{rom.MD5, game.Name}); err != nil {
				log.Fatal().Err(err).Msg("Failed to write to CSV file")
			}
		}
	}
	c.Flush()
	if err := c.Error(); err != nil {
		log.Fatal().Err(err).Msg("Failed to write to CSV file")
	}

	if err := f.Close(); err != nil {
		log.Fatal().Err(err).Msg("Failed to close CSV file")
	}
}
