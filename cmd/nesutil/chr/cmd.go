package chr

import (
	"gabe565.com/gones/cmd/nesutil/chr/decode"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chr",
		Short: "CHR graphics data utilities",
	}
	cmd.AddCommand(decode.New())
	return cmd
}
