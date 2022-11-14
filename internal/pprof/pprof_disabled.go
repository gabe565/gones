//go:build !pprof

package pprof

import "github.com/spf13/cobra"

func Flag(_ *cobra.Command) {}

func Spawn() {}
