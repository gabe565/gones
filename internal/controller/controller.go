package controller

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

type Controller struct {
	strobe bool
	index  byte
	bits   bitflags.Flags
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
