package test

import (
	"errors"
	"io"

	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/console"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
)

func stubConsole(r io.ReadSeeker) (*console.Console, error) {
	cart, err := cartridge.FromiNes(r)
	if err != nil {
		return nil, err
	}
	mapper, err := cartridge.NewMapper(cart)
	if err != nil {
		return nil, err
	}

	conf := config.NewDefault()
	c := &console.Console{
		Config:    conf,
		CPU:       nil,
		Bus:       nil,
		PPU:       ppu.New(conf, mapper),
		APU:       apu.New(conf),
		Cartridge: cart,
		Mapper:    mapper,
	}
	c.Bus = bus.New(conf, c.Mapper, c.PPU, c.APU)
	c.CPU = cpu.New(c.Bus)

	c.PPU.SetCPU(c.CPU)
	c.APU.SetCPU(c.CPU)

	return c, nil
}

type ConsoleTest struct {
	Console *console.Console
	ResetIn uint16

	Callback func(b *ConsoleTest) error
}

func NewConsoleTest(r io.ReadSeeker, callback func(console *ConsoleTest) error) (*ConsoleTest, error) {
	c, err := stubConsole(r)
	if err != nil {
		return nil, err
	}

	ct := &ConsoleTest{
		Console:  c,
		Callback: callback,
	}
	return ct, nil
}

func (b *ConsoleTest) Run() error {
	for {
		if b.ResetIn != 0 {
			b.ResetIn--
			if b.ResetIn == 0 {
				b.Console.Reset()
			}
		}

		if b.Callback != nil {
			if err := b.Callback(b); err != nil {
				if errors.Is(err, console.ErrExit) {
					return nil
				}
				return err
			}
		}

		if b.Console.Step(true); b.Console.CPU.StepErr != nil {
			return b.Console.CPU.StepErr
		}
	}
}
