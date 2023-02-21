package main

import (
	"github.com/gabe565/gones/cmd"
	_ "github.com/gabe565/gones/internal/pprof"
	log "github.com/sirupsen/logrus"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js public

var (
	Version = "next"
	Commit  = ""
)

func main() {
	if err := cmd.New(buildVersion()).Execute(); err != nil {
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
