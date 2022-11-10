package console

import (
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/joypad"
	"github.com/gabe565/gones/internal/ppu"
)

func New(path string, callback func(*ppu.PPU, *joypad.Joypad)) (*cpu.CPU, error) {
	cart, err := cartridge.FromiNes(path)
	if err != nil {
		return nil, err
	}

	b := bus.New(cart)
	b.Callback = callback
	return cpu.New(b), nil
}
