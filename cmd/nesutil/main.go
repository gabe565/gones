package main

import (
	"log/slog"
	"os"
	"strings"

	"gabe565.com/gones/cmd/nesutil/root"
	"gabe565.com/gones/cmd/options"
)

var version = ""

func main() {
	rootCmd := root.New(options.WithVersion(version))
	if err := rootCmd.Execute(); err != nil {
		for _, s := range strings.Split(err.Error(), "\n") {
			slog.Error(s)
		}
		os.Exit(1)
	}
}
