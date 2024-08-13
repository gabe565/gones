package gones

import (
	"image"
	"image/png"
	"io/fs"
	"log/slog"

	"github.com/gabe565/gones/assets"
)

func getWindowIcons() []image.Image {
	icons := make([]image.Image, 0, 3)

	if err := fs.WalkDir(assets.Icons, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		f, err := assets.Icons.Open(path)
		if err != nil {
			slog.Error("Failed to open icon", "error", err)
			return nil
		}
		defer func(f fs.File) {
			_ = f.Close()
		}(f)

		icon, err := png.Decode(f)
		if err != nil {
			slog.Error("Failed to decode icon", "error", err)
			return nil
		}

		icons = append(icons, icon)
		return nil
	}); err != nil {
		slog.Error("Failed to load icons", "error", err)
	}

	return icons
}
