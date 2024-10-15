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
	flag.StringP(FlagPRG, "p", "", "PRG ROM output file path")
	util.Must(cmd.MarkFlagRequired(FlagPRG))

	flag.StringP(FlagCHR, "c", "", "CHR ROM output file path")
	flag.Uint8P(FlagMapper, "m", 0, "INES mapper number")
	flag.StringP(FlagMirror, "n", "", "Type of nametable mirroring (one of horizontal, vertical, fourscreen)")
	flag.BoolP(FlagBattery, "b", false, "Enable battery/extra RAM")

	return cmd
}

var ErrUnknownMirror = errors.New("unknown mirror")

func run(cmd *cobra.Command, args []string) error {
	cart := cartridge.New()

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

	cart.Header.SetMapper(util.Must2(cmd.Flags().GetUint8(FlagMapper)))

	mirror := util.Must2(cmd.Flags().GetString(FlagMirror))
	switch strings.ToLower(mirror) {
	case "horizontal", "h":
		cart.Header.SetMirror(cartridge.Horizontal)
	case "vertical", "v":
		cart.Header.SetMirror(cartridge.Vertical)
	case "fourscreen", "f":
		cart.Header.SetMirror(cartridge.FourScreen)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownMirror, mirror)
	}

	cart.Header.SetBattery(util.Must2(cmd.Flags().GetBool(FlagBattery)))

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
