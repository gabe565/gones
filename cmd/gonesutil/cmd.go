package main

import (
	"os"

	"github.com/gabe565/gones/cmd/gonesutil/root"
)

func main() {
	rootCmd := root.New()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
