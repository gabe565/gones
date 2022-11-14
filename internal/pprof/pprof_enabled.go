//go:build pprof

package pprof

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
)

var address string

func Flag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&address, "pprof", "localhost:3000", "Enables pprof http listener")
}

func Spawn() {
	if address != "" {
		go func() {
			log.WithField("address", address).Info("starting pprof")
			if err := http.ListenAndServe(address, nil); err != nil {
				log.WithError(err).Error("failed to start pprof")
			}
		}()
	}
}
