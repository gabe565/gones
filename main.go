package main

import (
	"github.com/gabe565/gones/cmd"
	_ "net/http/pprof"
	"os"
)

func main() {
	if err := cmd.New(buildVersion()).Execute(); err != nil {
		os.Exit(1)
	}
}
