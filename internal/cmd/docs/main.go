package main

import (
	"os"

	"github.com/gabe565/gones/cmd/gones"
	gonesutil "github.com/gabe565/gones/cmd/gonesutil/root"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	output := "./docs"
	commands := []*cobra.Command{
		gones.New(),
		gonesutil.New(),
	}

	if err := os.RemoveAll(output); err != nil {
		log.Fatal().Err(err).Msg("Failed to remove existing dir")
	}

	if err := os.MkdirAll(output, 0o777); err != nil {
		log.Fatal().Err(err).Msg("Failed to create directory")
	}

	for _, cmd := range commands {
		if err := doc.GenMarkdownTree(cmd, output); err != nil {
			log.Fatal().Err(err).Msg("Failed to generate markdown")
		}
	}
}
