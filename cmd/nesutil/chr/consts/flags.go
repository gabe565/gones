package consts

import (
	"errors"
	"fmt"
	"image/color"
	"strings"

	"gabe565.com/utils/colorx"
	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	FlagPalette = "palette"

	TilesPerRow  = 16
	TileSize     = 8
	BytesPerTile = 16

	ImageWidth = TilesPerRow * TileSize
)

func PaletteFlag(fs *pflag.FlagSet) {
	fs.StringSliceP(FlagPalette, "p", []string{"000", "555", "AAA", "FFF"}, "Palette to use. Must contain 4 hex colors.")
}

var ErrInvalidPalette = errors.New("palette must contain 4 hex colors")

func LoadPalette(cmd *cobra.Command) (color.Palette, error) {
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
