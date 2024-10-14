package ines

import (
	"gabe565.com/gones/cmd/nesutil/ines/extract"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ines",
		Short: "Commands that work with INES files",
	}
	cmd.AddCommand(extract.New())
	return cmd
}
