package root

import (
	"gabe565.com/gones/cmd/nesutil/chr"
	"gabe565.com/gones/cmd/nesutil/genie"
	"gabe565.com/gones/cmd/nesutil/ines"
	"gabe565.com/gones/cmd/nesutil/ls"
	"gabe565.com/gones/cmd/options"
	"github.com/spf13/cobra"
)

func New(opts ...options.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nesutil",
		Short: "GoNES command-line utilities",

		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(ls.New(), ines.New(), chr.New(), genie.New())

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}
