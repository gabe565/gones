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

type consoleTest struct {
	console *console.Console
	resetIn int

	cb func(ct *consoleTest) error
}

func newConsoleTest(r io.ReadSeeker, cb func(c *consoleTest) error) (*consoleTest, error) {
	c, err := stubConsole(r)
	if err != nil {
		return nil, err
	}

	ct := &consoleTest{
		console: c,
		cb:      cb,
	}
	return ct, nil
}

func (c *consoleTest) run() error {
	for {
		if c.console.Step(true); c.console.CPU.StepErr != nil {
			return c.console.CPU.StepErr
		}

		if c.cb != nil {
			if err := c.cb(c); err != nil {
				if errors.Is(err, console.ErrExit) {
					return nil
				}
				return err
			}
		}

		if c.resetIn != 0 {
			c.resetIn--
			if c.resetIn == 0 {
				c.console.Reset()
			}
		}
	}
}
