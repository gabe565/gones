package console

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var ErrExit = errors.New("exit")

type Console struct {
	CPU *cpu.CPU
	Bus *bus.Bus
	PPU *ppu.PPU
	APU *apu.APU

	Cartridge *cartridge.Cartridge
	Mapper    cartridge.Mapper

	audioCtx      *audio.Context
	player        *audio.Player
	closeOnUpdate bool
	enableTrace   bool
	debug         Debug
}

func New(cart *cartridge.Cartridge) (*Console, error) {
	console := Console{Cartridge: cart}

	var err error
	console.Mapper, err = cartridge.NewMapper(cart)
	if err != nil {
		return &console, err
	}

	if cart.Battery {
		if err := console.LoadSram(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return &console, err
		}
	}

	console.PPU = ppu.New(console.Mapper)
	console.APU = apu.New()
	console.Bus = bus.New(console.Mapper, console.PPU, console.APU)
	console.CPU = cpu.New(console.Bus)

	console.Mapper.SetCpu(console.CPU)
	console.PPU.SetCpu(console.CPU)
	console.APU.SetCpu(console.CPU)

	console.CPU.Reset()

	if config.K.Bool("audio.enabled") {
		console.audioCtx = audio.NewContext(consts.AudioSampleRate)
		console.player, err = console.audioCtx.NewPlayer(console.APU)
		if err != nil {
			return &console, err
		}
		console.player.SetBufferSize(time.Second / 50)
		console.player.Play()
	} else {
		console.APU.Enabled = false
	}

	if config.K.Bool("state.resume") {
		if err := console.LoadState(0); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return &console, err
			}
		}
	}

	return &console, nil
}

func (c *Console) Close() error {
	if config.K.Bool("state.resume") {
		if err := c.SaveState(0); err != nil {
			return err
		}
	}
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
		c.Mapper.Step(
			c.PPU.Mask.BackgroundEnable || c.PPU.Mask.SpriteEnable,
			c.PPU.Scanline,
			c.PPU.Cycles,
		)
	}

	for i := uint(0); i < cycles; i += 1 {
		c.APU.Step()
	}

	return err
}

func (c *Console) Reset() {
	c.CPU.Reset()
	c.PPU.Reset()
	c.APU.Reset()
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
	img := c.PPU.Image()
	screen.WritePixels(img.Pix)
}

func (c *Console) CloseOnUpdate() {
	c.closeOnUpdate = true
}

func (c *Console) SetTrace(v bool) {
	c.enableTrace = v
}

func (c *Console) SetDebug(v bool) {
	if v {
		log.Info("Enable step debug")
		c.debug = DebugWait
	} else {
		c.debug = DebugDisabled
	}
}
