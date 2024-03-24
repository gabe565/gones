package controller

import (
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/controller/button"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player string

const (
	Player1 Player = "player1"
	Player2 Player = "player2"
)

func NewController(conf *config.Config, player Player) Controller {
	controller := Controller{
		Keymap:         NewKeymap(conf, player),
		turboDutyCycle: conf.Input.TurboDutyCycle,
	}
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

	turboDutyCycle uint16
	turbo          uint16
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
		j.index++
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
				j.buttons[button] = j.turbo < j.turboDutyCycle/2
				continue
			}
		}
		j.buttons[button] = pressed
	}

	// Directional safety
	if j.buttons[button.Left] && j.buttons[button.Right] {
		j.buttons[button.Right] = false
	}
	if j.buttons[button.Up] && j.buttons[button.Down] {
		j.buttons[button.Down] = false
	}

	if turboPressed {
		if j.turbo == j.turboDutyCycle-1 {
			j.turbo = 0
		} else {
			j.turbo++
		}
	} else {
		j.turbo = 0
	}
}
