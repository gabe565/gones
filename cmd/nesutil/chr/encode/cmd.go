package encode

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gabe565.com/gones/cmd/nesutil/chr/consts"
	"gabe565.com/gones/internal/util"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encode PNG [CHR data]",
		Short: "Encode a PNG file into NES CHR data",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  run,

		ValidArgsFunction: validArgs,
	}
	consts.PaletteFlag(cmd.Flags())
	return cmd
}

func validArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return []string{"png"}, cobra.ShellCompDirectiveFilterFileExt
	case 1:
		return util.CompleteROM(cmd, args, toComplete)
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

var (
	ErrInvalidImage = errors.New("invalid image")
	ErrImageWidth   = errors.New("image width must be " + strconv.Itoa(consts.ImageWidth))
)

func run(cmd *cobra.Command, args []string) error {
	palette, err := consts.LoadPalette(cmd)
	if err != nil {
		return err
	}

	input := args[0]
	imgReader, err := os.Open(input)
	if err != nil {
		return err
	}
	defer imgReader.Close()

	imgRaw, err := png.Decode(imgReader)
	if err != nil {
		return err
	}
	_ = imgReader.Close()

	img, ok := imgRaw.(subImager)
	if !ok {
		return ErrInvalidImage
	}

	if img.Bounds().Max.X != consts.ImageWidth {
		return fmt.Errorf("%w; got %d", ErrImageWidth, img.Bounds().Max.X)
	}

	size := img.Bounds().Max.Y * (consts.TilesPerRow * consts.BytesPerTile) / consts.TileSize
	chr := make([]byte, 0, size)
	for tile := range generateTiles(img) {
		for bitplane := range 2 {
			for y := range consts.TileSize {
				var chrByte byte
				for x := range consts.TileSize {
					color := tile.At(x+tile.Bounds().Min.X, y+tile.Bounds().Min.Y)
					colorIndex := palette.Index(color)
					colorIndex = colorIndex >> bitplane & 1
					chrByte |= byte(colorIndex) << (7 - x)
				}
				chr = append(chr, chrByte)
			}
		}
	}

	var output string
	if len(args) > 1 {
		output = args[1]
	} else {
		output = strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	}

	slog.Info("Writing CHR data", "path", output)
	return os.WriteFile(output, chr, 0o644)
}

type subImager interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

func generateTiles(img subImager) iter.Seq[image.Image] {
	return func(yield func(image.Image) bool) {
		bounds := img.Bounds().Max
		for y := range bounds.Y / consts.TileSize {
			for x := range bounds.X / consts.TileSize {
				rect := image.Rect(x*consts.TileSize, y*consts.TileSize, (x+1)*consts.TileSize, (y+1)*consts.TileSize)
				subImg := img.SubImage(rect)
				if !yield(subImg) {
					return
				}
			}
		}
	}
}
