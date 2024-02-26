package palette

import (
	"encoding/binary"
	"image/color"
	"io"
)

type PalColor struct {
	R, G, B byte
}

func (l PalColor) RGBA() color.RGBA {
	return color.RGBA{
		R: l.R,
		G: l.G,
		B: l.B,
		A: 0xFF,
	}
}

func LoadPal(r io.Reader) error {
	var c PalColor
	for i := 0; i < len(Default); i++ {
		if err := binary.Read(r, binary.LittleEndian, &c); err != nil {
			return err
		}
		Default[i] = c.RGBA()
	}
	UpdateEmphasized()
	return nil
}
