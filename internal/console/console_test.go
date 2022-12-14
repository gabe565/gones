package console

import (
	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
	"io"
)

func stubConsole(r io.ReadSeeker) (*Console, error) {
	cart, err := cartridge.FromiNes(r)
	if err != nil {
		return nil, err
	}
	mapper, err := cartridge.NewMapper(cart)
	if err != nil {
		return nil, err
	}
	console := Console{Cartridge: cart, Mapper: mapper}

	console.PPU = ppu.New(console.Mapper)
	console.APU = apu.New()
	console.Bus = bus.New(console.Mapper, console.PPU, console.APU)
	console.CPU = cpu.New(console.Bus)

	console.Mapper.SetCpu(console.CPU)
	console.PPU.SetCpu(console.CPU)
	console.APU.SetCpu(console.CPU)

	return &console, nil
}
