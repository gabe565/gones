package main

import (
	"log/slog"
	"os"

	"gabe565.com/gones/cmd/gones"
	"gabe565.com/gones/cmd/options"
	_ "gabe565.com/gones/internal/pprof"
	"github.com/spf13/cobra"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js web/src/scripts
//go:generate sh -c "gzip -c internal/database/database.csv > internal/database/database.csv.gz"

var version = "beta"

func main() {
	cobra.MousetrapHelpText = ""
	rootCmd := gones.New(options.WithVersion(version))
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
