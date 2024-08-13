package main

import (
	"log/slog"

	"github.com/gabe565/gones/cmd/gones"
	"github.com/gabe565/gones/cmd/options"
	_ "github.com/gabe565/gones/internal/pprof"
	"github.com/spf13/cobra"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js web/src/scripts
//go:generate sh -c "gzip -c internal/database/database.csv > internal/database/database.csv.gz"

var version = "beta"

func main() {
	cobra.MousetrapHelpText = ""
	rootCmd := gones.New(options.WithVersion(version))
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Exiting due to an error", "error", err)
	}
}
