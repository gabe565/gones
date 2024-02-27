package main

import (
	"os"

	"github.com/gabe565/gones/internal/config"
	"github.com/pelletier/go-toml/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	f, err := os.Create("config_example.toml")
	if err != nil {
		log.Panic(err)
	}

	encoder := toml.NewEncoder(f)
	conf := config.NewDefault()
	if err := encoder.Encode(conf); err != nil {
		log.Panic(err)
	}
}
