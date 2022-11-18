package controller

import (
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ButtonA = 1 << iota
	ButtonB
	Select
	Start
	Up
	Down
	Left
	Right
)

type Controller struct {
	strobe bool
	index  byte
	bits   bitflags.Flags

	Keymap Keymap[ebiten.Key]
	turbo  uint8
}

func (j *Controller) Write(data byte) {
	j.strobe = data&1 == 1
	if j.strobe {
		j.index = 0
	}
}

func (j *Controller) Read() byte {
	if j.index >= 8 {
		return 1
	}

	var value byte
	if j.bits.Intersects(1 << j.index) {
		value = 1
	}

	if !j.strobe && j.index < 8 {
		j.index += 1
	}
	return value
}

func (j *Controller) UpdateInput() {
	for key, button := range j.Keymap.Regular {
		keyPressed := ebiten.IsKeyPressed(key)
		j.bits.Set(button, keyPressed)
	}

	var turboPressed bool
	for key, button := range j.Keymap.Turbo {
		if ebiten.IsKeyPressed(key) && !j.bits.Intersects(button) {
			turboPressed = true
			j.bits.Set(button, j.turbo%6 < 3)
		}
	}
	if turboPressed {
		j.turbo += 1
	} else {
		j.turbo = 0
	}
}
