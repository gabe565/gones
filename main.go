package main

import (
	"github.com/gabe565/gones/cmd"
	_ "github.com/gabe565/gones/internal/pprof"
	log "github.com/sirupsen/logrus"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js public

func main() {
	if err := cmd.New(buildVersion()).Execute(); err != nil {
		log.Fatal(err)
	}
}
