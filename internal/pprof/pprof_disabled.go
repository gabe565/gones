//go:build !pprof

package pprof

const Enabled = false

func ListenAndServe() error {
	return nil
}
