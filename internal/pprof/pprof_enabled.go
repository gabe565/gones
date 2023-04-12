//go:build pprof

package pprof

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
)

var address = "localhost:3000"

func init() {
	go func() {
		log.WithField("address", address).Info("starting pprof")
		if err := http.ListenAndServe(address, nil); err != nil {
			log.WithError(err).Error("failed to start pprof")
		}
	}()
}
