package console

import (
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
)

type Console struct {
	CPU *cpu.CPU
	Bus *bus.Bus
	PPU *ppu.PPU
}

func New(path string) (Console, error) {
	var console Console

	cart, err := cartridge.FromiNesFile(path)
	if err != nil {
		return console, err
	}

	console.PPU = ppu.New(cart)
	console.Bus = bus.New(cart, console.PPU)
	console.CPU = cpu.New(console.Bus)

	console.CPU.Reset()

	return console, nil
}
