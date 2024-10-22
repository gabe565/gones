package extract

import (
	"bytes"
	"encoding/binary"
	"log/slog"
	"os"

	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/util"
	"github.com/spf13/cobra"
)

const (
	FlagHeader = "header"
	FlagPRG    = "prg"
	FlagCHR    = "chr"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract ROM",
		Short: "Extract PRG/CHR ROM data from an INES ROM",
		Args:  cobra.ExactArgs(1),
		RunE:  run,

		ValidArgsFunction: util.CompleteROM,
	}

	flag := cmd.Flags()
	flag.StringP(FlagHeader, "H", "", "Header output file path")
	flag.StringP(FlagPRG, "p", "", "PRG ROM output file path")
	flag.StringP(FlagCHR, "c", "", "CHR ROM output file path")
	cmd.MarkFlagsOneRequired(FlagPRG, FlagCHR)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	path := args[0]
	cart, err := cartridge.FromINESFile(path)
	if err != nil {
		return err
	}

	if header := util.Must2(cmd.Flags().GetString(FlagHeader)); header != "" {
		slog.Info("Extracting header", "path", header)
		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.LittleEndian, cart.Header); err != nil {
			return err
		}
		if err := os.WriteFile(header, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}

	if prg := util.Must2(cmd.Flags().GetString(FlagPRG)); prg != "" {
		slog.Info("Extracting PRG", "path", prg)
		if err := os.WriteFile(prg, cart.PRG, 0o644); err != nil {
			return err
		}
	}

	if chr := util.Must2(cmd.Flags().GetString(FlagCHR)); chr != "" {
		if cart.Header.CHRCount == 0 {
			slog.Warn("Game does not have CHR. Skipping")
		} else {
			slog.Info("Extracting CHR", "path", chr)
			if err := os.WriteFile(chr, cart.CHR, 0o644); err != nil {
				return err
			}
		}
	}

	return nil
}
