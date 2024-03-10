package main

import (
	"github.com/gabe565/gones/cmd/gones"
	_ "github.com/gabe565/gones/internal/pprof"
	log "github.com/sirupsen/logrus"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js web/src/scripts

const (
	Version = "next"
	Commit  = ""
)

func main() {
	rootCmd := gones.New()
	rootCmd.Version = buildVersion()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func buildVersion() string {
	result := Version
	if Commit != "" {
		result += " (" + Commit + ")"
	}
	return result
}
