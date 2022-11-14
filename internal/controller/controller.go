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

	Keymap    Keymap
	GamepadID ebiten.GamepadID
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
	if j.bits.Has(1 << j.index) {
		value = 1
	}

	if !j.strobe && j.index < 8 {
		j.index += 1
	}
	return value
}

func (j *Controller) Set(button bitflags.Flags, status bool) {
	j.bits.Set(button, status)
}

func (j *Controller) UpdateInput() {
	for key, button := range j.Keymap {
		keyPressed := ebiten.IsKeyPressed(key)
		gamepadPressed := ebiten.IsGamepadButtonPressed(
			j.GamepadID,
			ebiten.GamepadButton(Joystick[button]),
		)
		j.Set(button, keyPressed || gamepadPressed)
	}
}
