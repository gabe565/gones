package main

import (
	"encoding/csv"
	"github.com/gabe565/gones/internal/database/nointro"
	log "github.com/sirupsen/logrus"
	"os"
)

//go:generate go run .

func main() {
	datafile, err := nointro.Load(nointro.Nes)
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Create("../database.csv")
	if err != nil {
		log.Panic(err)
	}

	c := csv.NewWriter(f)
	for _, game := range datafile.Games {
		for _, rom := range game.Roms {
			if err := c.Write([]string{rom.MD5, game.Name}); err != nil {
				log.Panic(err)
			}
		}
	}
	c.Flush()
	if err := c.Error(); err != nil {
		log.Panic(err)
	}

	if err := f.Close(); err != nil {
		log.Panic(err)
	}
}
