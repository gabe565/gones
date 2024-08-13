//go:build pprof

package pprof

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
)

var address = "localhost:3000"

func init() { //nolint:all
	go func() {
		slog.Info("Starting pprof", "address", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			slog.Error("Failed to start pprof", "error", err)
		}
	}()
}
