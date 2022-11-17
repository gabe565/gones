package console

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"os"
)

var ErrExit = errors.New("exit")

type Console struct {
	CPU *cpu.CPU
	Bus *bus.Bus
	PPU *ppu.PPU

	Cartridge *cartridge.Cartridge

	closeOnUpdate bool
	enableTrace   bool
	debug         Debug
}

func New(path string) (*Console, error) {
	var console Console

	cart, err := cartridge.FromiNesFile(path)
	if err != nil {
		return &console, err
	}

	console.Cartridge = cart
	if cart.Battery {
		if err := console.LoadSram(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return &console, err
		}
	}

	console.PPU = ppu.New(cart)
	console.Bus = bus.New(cart, console.PPU)
	console.CPU = cpu.New(console.Bus)

	console.CPU.Reset()

	return &console, nil
}

func (c *Console) Close() error {
	if c.Cartridge.Battery {
		if err := c.SaveSram(); err != nil {
			return err
		}
	}
	return nil
}

var ErrRender = errors.New("render triggered")

func (c *Console) Step() error {
	if c.enableTrace {
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

func (c *Console) Layout(_, _ int) (int, int) {
	return ppu.Width, ppu.TrimmedHeight
}

func (c *Console) Update() error {
	if ebiten.IsWindowBeingClosed() || c.closeOnUpdate {
		if err := c.Close(); err != nil {
			return err
		}
		return ErrExit
	}

	c.CheckInput()

	if c.debug == DebugWait {
		return nil
	}

	for {
		if err := c.Step(); err != nil {
			if errors.Is(err, ErrRender) {
				break
			}
			return err
		}

		if c.debug == DebugStepFrame {
			break
		}
	}

	if c.debug != DebugDisabled {
		c.debug = DebugWait
	}

	return nil
}

func (c *Console) Draw(screen *ebiten.Image) {
	img := ebiten.NewImageFromImage(c.PPU.Render())
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(0, -ppu.TrimHeight)
	screen.DrawImage(img, &op)
}

func (c *Console) CloseOnUpdate() {
	c.closeOnUpdate = true
}
