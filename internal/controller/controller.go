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
	turbo  bool
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
	var turboPressed bool
	for button, key := range j.Keymap.Regular {
		pressed := ebiten.IsKeyPressed(key)
		if !pressed {
			turboKey, ok := j.Keymap.Turbo[button]
			if ok && ebiten.IsKeyPressed(turboKey) {
				turboPressed = true
				j.buttons[button] = j.turbo
				continue
			}
		}
		j.buttons[button] = pressed
	}
	if turboPressed {
		j.turbo = !j.turbo
	} else {
		j.turbo = false
	}
}
