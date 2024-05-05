//go:build freebsd || linux || netbsd || openbsd || windows

package gones

import "github.com/hajimehoshi/ebiten/v2"

func init() { //nolint:all
	go func() {
		ebiten.SetWindowIcon(getWindowIcons())
	}()
}
