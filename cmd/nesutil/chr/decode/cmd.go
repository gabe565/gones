package decode

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gabe565.com/gones/cmd/nesutil/chr/consts"
	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/util"
	"gabe565.com/utils/colorx"
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

const FlagPalette = "palette"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode { ROM | CHR data } [output]",
		Short: "Decode NES CHR data into a PNG file",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  run,

		ValidArgsFunction: util.CompleteROM,
	}
	fs := cmd.Flags()
	fs.StringSliceP(FlagPalette, "p", []string{"000", "555", "AAA", "FFF"}, "Palette to use. Must contain 4 hex colors.")
	return cmd
}

var ErrNoCHR = errors.New("ROM file has no CHR data")

func run(cmd *cobra.Command, args []string) error {
	palette, err := loadPalette(cmd)
	if err != nil {
		return err
	}

	input := args[0]
	chr, err := loadCHR(input)
	if err != nil {
		return err
	}

	height := len(chr) / (consts.TilesPerRow * consts.BytesPerTile) * consts.TileSize
	img := image.NewPaletted(image.Rect(0, 0, consts.TilesPerRow*consts.TileSize, height), palette)

	var count int
	for i, tile := range generateTiles(chr) {
		xOffset := (i % consts.TilesPerRow) * consts.TileSize
		yOffset := (i / consts.TilesPerRow) * consts.TileSize

		for y := range consts.TileSize {
			for x := range consts.TileSize {
				img.SetColorIndex(xOffset+x, yOffset+y, tile[y*consts.TileSize+x])
			}
		}
		count++
	}

	var output string
	if len(args) > 1 {
		output = args[1]
	} else {
		output = strings.TrimSuffix(filepath.Base(input), filepath.Ext(input)) + ".png"
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	slog.Info("Wrote file", "path", output, "tiles", count)
	return nil
}

var ErrInvalidPalette = errors.New("palette must contain 4 hex colors")

func loadPalette(cmd *cobra.Command) (color.Palette, error) {
	paletteRaw := must.Must2(cmd.Flags().GetStringSlice(FlagPalette))
	if len(paletteRaw) != 4 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPalette, strings.Join(paletteRaw, ","))
	}

	palette := make(color.Palette, 0, len(paletteRaw))
	for _, str := range paletteRaw {
		c, err := colorx.ParseHex(str)
		if err != nil {
			return nil, err
		}
		palette = append(palette, c)
	}

	return palette, nil
}

func loadCHR(input string) ([]byte, error) {
	if filepath.Ext(input) == ".nes" {
		cart, err := cartridge.FromINESFile(input)
		if err != nil {
			return nil, err
		}

		if cart.Header.CHRCount == 0 {
			return nil, fmt.Errorf("%w: %s", ErrNoCHR, input)
		}

		return cart.CHR, nil
	}

	return os.ReadFile(input)
}

func generateTiles(chr []byte) iter.Seq2[int, []byte] {
	return func(yield func(int, []byte) bool) {
		count := len(chr) / consts.BytesPerTile

		for i := range count {
			tile := chr[i*consts.BytesPerTile : (i+1)*consts.BytesPerTile]
			decodedTile := make([]byte, 0, consts.TileSize*consts.TileSize)

			for y := range consts.TileSize {
				for s := range consts.TileSize {
					x := 7 - s
					lo := (tile[y] >> x) & 1
					hi := (tile[y+consts.TileSize] >> x) & 1
					decodedTile = append(decodedTile, hi<<1|lo)
				}
			}

			if !yield(i, decodedTile) {
				return
			}
		}
	}
}
