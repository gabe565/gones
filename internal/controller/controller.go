package controller

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Button uint8

const (
	ButtonA Button = iota
	ButtonB
	Select
	Start
	Up
	Down
	Left
	Right
)

type Controller struct {
	Enabled bool
	strobe  bool
	index   byte
	buttons [8]bool

	Keymap Keymap
	turbo  uint8
}

func (j *Controller) Write(data byte) {
	j.strobe = data&1 == 1
	if j.strobe {
		j.index = 0
	}
}

func (j *Controller) Read() byte {
	if !j.Enabled {
		return 0
	}

	if j.index >= 8 {
		return 1
	}

	var value byte
	if j.buttons[j.index] {
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
		j.buttons[button] = keyPressed
	}

	var turboPressed bool
	for key, button := range j.Keymap.Turbo {
		if ebiten.IsKeyPressed(key) && !j.buttons[button] {
			turboPressed = true
			j.buttons[button] = j.turbo%6 < 3
		}
	}
	if turboPressed {
		j.turbo += 1
	} else {
		j.turbo = 0
	}
}
