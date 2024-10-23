package decode

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gabe565.com/gones/cmd/nesutil/chr/consts"
	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/util"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode { ROM | CHR data } [output]",
		Short: "Decode NES CHR data into a PNG file",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  run,

		ValidArgsFunction: util.CompleteROM,
	}
	consts.PaletteFlag(cmd.Flags())
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	palette, err := consts.LoadPalette(cmd)
	if err != nil {
		return err
	}

	input := args[0]
	chr, err := loadCHR(input)
	if err != nil {
		return err
	}

	height := len(chr) / (consts.TilesPerRow * consts.BytesPerTile) * consts.TileSize
	img := image.NewPaletted(image.Rect(0, 0, consts.ImageWidth, height), palette)

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

var ErrNoCHR = errors.New("ROM file has no CHR data")

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
