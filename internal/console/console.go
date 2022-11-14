package console

import (
	"errors"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
)

type Console struct {
	CPU *cpu.CPU
	Bus *bus.Bus
	PPU *ppu.PPU

	EnableTrace bool
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

var ErrRender = errors.New("render triggered")

func (c *Console) Step() error {
	if c.EnableTrace {
		fmt.Println(c.CPU.Trace())
	}

	cycles, err := c.CPU.Step()
	if err != nil {
		return err
	}

	for i := uint(0); i < cycles*3; i += 1 {
		if c.PPU.Step() {
			err = ErrRender
		}
	}

	select {
	case interrupt := <-c.PPU.Interrupt():
		c.CPU.Interrupt <- interrupt
	default:
		//
	}

	return err
}

func (c *Console) Reset() {
	c.CPU.Reset()
	c.PPU.Reset()
}

func (c *Console) UpdateInput(win *pixelgl.Window) {
	c.Bus.Controller1.UpdateInput(win)
	c.Bus.Controller2.UpdateInput(win)
}
