package main

import (
	"github.com/gabe565/gones/cmd/gones"
	_ "github.com/gabe565/gones/internal/pprof"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js web/src/scripts
//go:generate sh -c "gzip -c internal/database/database.csv > internal/database/database.csv.gz"

const (
	Version = "next"
	Commit  = ""
)

func main() {
	cobra.MousetrapHelpText = ""
	rootCmd := gones.New()
	rootCmd.Version = buildVersion()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Exiting due to an error")
	}
}

func buildVersion() string {
	result := Version
	if Commit != "" {
		result += " (" + Commit + ")"
	}
	return result
}
