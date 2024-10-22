package create

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/util"
	"github.com/spf13/cobra"
)

const (
	FlagHeader  = "header"
	FlagPRG     = "prg"
	FlagCHR     = "chr"
	FlagMapper  = "mapper"
	FlagMirror  = "mirror"
	FlagBattery = "battery"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create ROM",
		Short: "Create an INES ROM file",
		Args:  cobra.ExactArgs(1),
		RunE:  run,

		ValidArgsFunction: util.CompleteROM,
	}

	flag := cmd.Flags()
	flag.StringP(FlagHeader, "H", "", "Header file")
	flag.StringP(FlagPRG, "p", "", "PRG ROM output file path")
	flag.StringP(FlagCHR, "c", "", "CHR ROM output file path")
	flag.Uint8P(FlagMapper, "m", 0, "INES mapper number")
	flag.StringP(FlagMirror, "n", "", "Type of nametable mirroring (one of horizontal, vertical, fourscreen)")
	flag.BoolP(FlagBattery, "b", false, "Enable battery/extra RAM")
	util.Must(cmd.MarkFlagRequired(FlagPRG))

	return cmd
}

var ErrUnknownMirror = errors.New("unknown mirror")

func run(cmd *cobra.Command, args []string) error {
	cart := cartridge.New()

	if header := util.Must2(cmd.Flags().GetString(FlagHeader)); header != "" {
		slog.Info("Loading header", "path", header)
		f, err := os.Open(header)
		if err != nil {
			return err
		}

		if err := binary.Read(f, binary.LittleEndian, &cart.Header); err != nil {
			return err
		}

		_ = f.Close()
	}

	if prg := util.Must2(cmd.Flags().GetString(FlagPRG)); prg != "" {
		slog.Info("Loading PRG", "path", prg)
		var err error
		if cart.PRG, err = os.ReadFile(prg); err != nil {
			return err
		}
		cart.Header.PRGCount = byte(len(cart.PRG) / consts.PRGChunkSize)
	}

	if chr := util.Must2(cmd.Flags().GetString(FlagCHR)); chr != "" {
		slog.Info("Loading CHR", "path", chr)
		var err error
		if cart.CHR, err = os.ReadFile(chr); err != nil {
			return err
		}
		cart.Header.CHRCount = byte(len(cart.CHR) / consts.CHRChunkSize)
	}

	if cmd.Flags().Lookup(FlagMapper).Changed {
		mapper := util.Must2(cmd.Flags().GetUint8(FlagMapper))
		slog.Info("Set mapper", "value", mapper)
		cart.Header.SetMapper(mapper)
	}

	if cmd.Flags().Lookup(FlagMirror).Changed {
		var mirror cartridge.Mirror
		switch strings.ToLower(util.Must2(cmd.Flags().GetString(FlagMirror))) {
		case "horizontal", "h":
			mirror = cartridge.Horizontal
		case "vertical", "v":
			mirror = cartridge.Vertical
		case "fourscreen", "f":
			mirror = cartridge.FourScreen
		default:
			return fmt.Errorf("%w: %s", ErrUnknownMirror, mirror)
		}
		slog.Info("Set mirror", "value", mirror)
		cart.Header.SetMirror(mirror)
	}

	if cmd.Flags().Lookup(FlagBattery).Changed {
		battery := util.Must2(cmd.Flags().GetBool(FlagBattery))
		slog.Info("Set battery", "value", battery)
		cart.Header.SetBattery(battery)
	}

	f, err := os.Create(args[0])
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	// Write header
	if err := binary.Write(f, binary.LittleEndian, cart.Header); err != nil {
		return err
	}

	// Write PRG
	if _, err := f.Write(cart.PRG); err != nil {
		return err
	}

	if len(cart.CHR) != 0 {
		// Write CHR
		if _, err := f.Write(cart.CHR); err != nil {
			return err
		}
	}

	return f.Close()
}
