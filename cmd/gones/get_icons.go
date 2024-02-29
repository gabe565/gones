package gones

import (
	"image"
	"io/fs"

	"github.com/gabe565/gones/assets"
	log "github.com/sirupsen/logrus"
)

func getWindowIcons() []image.Image {
	icons := make([]image.Image, 0, 3)

	err := fs.WalkDir(assets.Icons, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		f, err := assets.Icons.Open(path)
		if err != nil {
			log.WithError(err).Error("Failed to open icon")
			return nil
		}
		defer func(f fs.File) {
			_ = f.Close()
		}(f)

		icon, _, err := image.Decode(f)
		if err != nil {
			log.WithError(err).Error("Failed to decode icon")
			return nil
		}

		icons = append(icons, icon)
		return nil
	})
	if err != nil {
		log.WithError(err).Error("Failed to load icons")
	}

	return icons
}
