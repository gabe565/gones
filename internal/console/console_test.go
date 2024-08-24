package console

import (
	"io"

	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
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

	conf := config.NewDefault()
	console := Console{
		config:    &conf,
		Cartridge: cart,
		Mapper:    mapper,
	}

	console.PPU = ppu.New(config.Overscan{}, console.Mapper)
	console.APU = apu.New(&conf)
	console.Bus = bus.New(&conf, console.Mapper, console.PPU, console.APU)
	console.CPU = cpu.New(console.Bus)

	console.PPU.SetCPU(console.CPU)
	console.APU.SetCPU(console.CPU)

	return &console, nil
}
