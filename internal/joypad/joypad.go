package joypad

import "github.com/gabe565/gones/internal/bitflags"

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

type Joypad struct {
	strobe       bool
	buttonIdx    byte
	buttonStatus bitflags.Flags
}

func (j *Joypad) Write(data byte) {
	j.strobe = data&1 == 1
	if j.strobe {
		j.buttonIdx = 0
	}
}

func (j *Joypad) Read() byte {
	if j.buttonIdx > 7 {
		return 1
	}

	status := byte(j.buttonStatus)
	response := (status & (1 << j.buttonIdx)) >> j.buttonIdx
	if !j.strobe && j.buttonIdx <= 7 {
		j.buttonIdx += 1
	}
	return response
}

func (j *Joypad) Set(button bitflags.Flags, status bool) {
	j.buttonStatus.Set(button, status)
}
