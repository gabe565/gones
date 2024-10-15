package ines

import (
	"gabe565.com/gones/cmd/nesutil/ines/create"
	"gabe565.com/gones/cmd/nesutil/ines/extract"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ines",
		Short: "INES ROM utilities",
	}
	cmd.AddCommand(extract.New(), create.New())
	return cmd
}
