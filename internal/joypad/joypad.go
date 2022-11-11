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
	strobe bool
	index  byte
	bits   bitflags.Flags
}

func (j *Joypad) Write(data byte) {
	j.strobe = data&1 == 1
	if j.strobe {
		j.index = 0
	}
}

func (j *Joypad) Read() byte {
	if j.index > 7 {
		return 1
	}

	response := (byte(j.bits) & (1 << j.index)) >> j.index
	if !j.strobe && j.index <= 7 {
		j.index += 1
	}
	return response
}

func (j *Joypad) Set(button bitflags.Flags, status bool) {
	j.bits.Set(button, status)
}
