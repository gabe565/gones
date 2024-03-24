package palette

import (
	"encoding/binary"
	"image/color"
	"io"
	"os"
	"path/filepath"

	"github.com/gabe565/gones/internal/config"
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
	for i := range len(Default) {
		if err := binary.Read(r, binary.LittleEndian, &c); err != nil {
			return err
		}
		Default[i] = c.RGBA()
	}
	UpdateEmphasized()
	return nil
}

func LoadPalFile(path string) error {
	if path == "" {
		UpdateEmphasized()
		return nil
	}

	if !filepath.IsAbs(path) {
		palDir, err := config.GetPaletteDir()
		if err != nil {
			return err
		}

		path = filepath.Join(palDir, path)
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	return LoadPal(f)
}
