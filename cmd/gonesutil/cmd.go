package main

import (
	"os"

	"github.com/gabe565/gones/cmd/gonesutil/root"
	"github.com/gabe565/gones/cmd/options"
)

var version = ""

func main() {
	rootCmd := root.New(options.WithVersion(version))
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
