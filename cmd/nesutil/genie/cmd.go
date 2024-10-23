package genie

import (
	"gabe565.com/gones/cmd/nesutil/genie/decode"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genie",
		Short: "Game Genie code utilities",
	}
	cmd.AddCommand(decode.New())
	return cmd
}
