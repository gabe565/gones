package main

import (
	"log/slog"
	"os"

	"gabe565.com/gones/cmd/nesutil/root"
	"gabe565.com/gones/cmd/options"
)

var version = ""

func main() {
	rootCmd := root.New(options.WithVersion(version))
	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
