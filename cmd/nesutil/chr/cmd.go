package chr

import (
	"gabe565.com/gones/cmd/nesutil/chr/decode"
	"gabe565.com/gones/cmd/nesutil/chr/encode"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chr",
		Short: "CHR graphics data utilities",
	}
	cmd.AddCommand(encode.New(), decode.New())
	return cmd
}
