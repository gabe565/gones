//go:build !js

package gones

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gabe565/gones/cmd/options"
	"github.com/gabe565/gones/internal/config"
	"github.com/spf13/cobra"
)

func New(opts ...options.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gones ROM",
		Short: "NES emulator written in Go",
		RunE:  runCobra,
		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"nes"}, cobra.ShellCompDirectiveFilterFileExt
		},

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
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}
	cmd.SilenceUsage = true

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	return run(ctx, conf, path)
}
