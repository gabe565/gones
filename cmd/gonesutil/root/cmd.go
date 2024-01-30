package root

import (
	"github.com/gabe565/gones/cmd/gonesutil/ls"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gonesutil",
		Short: "GoNES command-line utilities",

		DisableAutoGenTag: true,
	}

	cmd.AddCommand(ls.New())

	return cmd
}
