package controller

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player string

const (
	Player1 Player = "player1"
	Player2        = "player2"
)

func NewController(player Player) Controller {
	controller := Controller{Keymap: NewKeymap(player)}
	if len(controller.Keymap.Regular) != 0 {
		controller.Enabled = true
	}
	return controller
}

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

	// Directional safety
	if j.buttons[Left] && j.buttons[Right] {
		j.buttons[Right] = false
	}
	if j.buttons[Up] && j.buttons[Down] {
		j.buttons[Down] = false
	}

	if turboPressed {
		j.turbo = !j.turbo
	} else {
		j.turbo = false
	}
}
