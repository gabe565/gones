//go:build pprof

package pprof

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/rs/zerolog/log"
)

var address = "localhost:3000"

func init() { //nolint:all
	go func() {
		log.Info().Str("address", address).Msg("starting pprof")
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Err(err).Msg("Failed to start pprof")
		}
	}()
}
