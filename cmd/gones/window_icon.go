//go:build freebsd || linux || netbsd || openbsd || windows

package gones

import "github.com/hajimehoshi/ebiten/v2"

//nolint:gochecknoinits
func init() {
	go func() {
		ebiten.SetWindowIcon(getWindowIcons())
	}()
}
