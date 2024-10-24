//go:build !js

package gones

import (
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/gones/cmd/options"
	"gabe565.com/gones/internal/config"
	"gabe565.com/gones/internal/util"
	"github.com/spf13/cobra"
)

func New(opts ...options.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gones ROM",
		Short: "NES emulator written in Go",
		RunE:  runCobra,

		ValidArgsFunction: util.CompleteROM,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}
	config.Flags(cmd)

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func runCobra(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	cart, err := loadCartridge(path)
	if err != nil {
		return err
	}

	conf := config.NewDefault()
	if err := conf.Load(cmd, cart.Name(), cart.Hash()); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	return run(ctx, conf, cart)
}
