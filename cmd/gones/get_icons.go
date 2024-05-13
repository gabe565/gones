package gones

import (
	"image"
	"image/png"
	"io/fs"

	"github.com/gabe565/gones/assets"
	"github.com/rs/zerolog/log"
)

func getWindowIcons() []image.Image {
	icons := make([]image.Image, 0, 3)

	if err := fs.WalkDir(assets.Icons, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		f, err := assets.Icons.Open(path)
		if err != nil {
			log.Err(err).Msg("Failed to open icon")
			return nil
		}
		defer func(f fs.File) {
			_ = f.Close()
		}(f)

		icon, err := png.Decode(f)
		if err != nil {
			log.Err(err).Msg("Failed to decode icon")
			return nil
		}

		icons = append(icons, icon)
		return nil
	}); err != nil {
		log.Err(err).Msg("Failed to load icons")
	}

	return icons
}
