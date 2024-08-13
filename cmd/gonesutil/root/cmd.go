package root

import (
	"github.com/gabe565/gones/cmd/gonesutil/ls"
	"github.com/gabe565/gones/cmd/options"
	"github.com/spf13/cobra"
)

func New(opts ...options.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gonesutil",
		Short: "GoNES command-line utilities",

		DisableAutoGenTag: true,
	}
	cmd.AddCommand(ls.New())

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}
