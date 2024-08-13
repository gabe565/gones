//go:build pprof

package pprof

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const (
	Enabled = true
	address = "localhost:3000"
)

func ListenAndServe() error {
	slog.Info("Starting pprof", "address", address)
	server := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}
