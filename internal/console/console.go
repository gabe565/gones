package console

import (
	"errors"
	"fmt"
	"os"
	"time"

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
)

const AutoSaveNum = 0

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

	autosave *time.Ticker
	rate     uint8
}

func New(cart *cartridge.Cartridge) (*Console, error) {
	console := Console{Cartridge: cart, rate: 1}

	var err error
	console.Mapper, err = cartridge.NewMapper(cart)
	if err != nil {
		return &console, err
	}

	if err := console.LoadSram(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return &console, err
	}

	console.PPU = ppu.New(console.Mapper)
	console.APU = apu.New()
	console.Bus = bus.New(console.Mapper, console.PPU, console.APU)
	console.CPU = cpu.New(console.Bus)

	if mapper, ok := console.Mapper.(cartridge.MapperInterrupts); ok {
		mapper.SetCpu(console.CPU)
	}
	console.PPU.SetCpu(console.CPU)
	console.APU.SetCpu(console.CPU)

	if config.K.Bool("audio.enabled") {
		console.audioCtx = audio.NewContext(consts.AudioSampleRate)
		console.player, err = console.audioCtx.NewPlayer(console.APU)
		if err != nil {
			return &console, err
		}
		console.player.SetBufferSize(time.Second / 20)
		console.player.Play()
	} else {
		console.APU.Enabled = false
	}

	if config.K.Bool("state.resume") {
		if err := console.LoadState(AutoSaveNum); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return &console, err
			}
		}
	}

	if duration := config.K.Duration("state.interval"); duration != 0 {
		console.autosave = time.NewTicker(config.K.Duration("state.interval"))
	}

	return &console, nil
}

func (c *Console) Close() error {
	c.autosave.Stop()
	if config.K.Bool("state.resume") {
		if err := c.SaveState(AutoSaveNum); err != nil {
			return err
		}
	}
	if err := c.SaveSram(); err != nil {
		return err
	}
	return nil
}

func (c *Console) Step(render bool) error {
	if c.enableTrace {
		fmt.Println(c.Trace())
	}

	cycles, err := c.CPU.Step()
	if err != nil {
		return err
	}
	if mapper, ok := c.Mapper.(cartridge.MapperOnCPUStep); ok {
		mapper.OnCPUStep()
	}

	for i := uint(0); i < cycles*3; i += 1 {
		c.PPU.Step(render)
	}

	for i := uint(0); i < cycles; i += 1 {
		c.APU.Step()
	}

	if c.autosave != nil {
		select {
		case <-c.autosave.C:
			if err := c.SaveSram(); err != nil {
				log.WithError(err).Error("Auto-save failed")
			}
			if err := c.SaveState(AutoSaveNum); err != nil {
				log.WithError(err).Error("State auto-save failed")
			}
		default:
		}
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
	if c.closeOnUpdate {
		return ErrExit
	}

	c.CheckInput()

	if c.debug == DebugWait {
		return nil
	}

	for i := uint8(0); i < c.rate; i += 1 {
		if c.rate != 1 {
			c.PPU.RenderDone = false
		}
		for {
			if err := c.Step(i == c.rate-1); err != nil {
				return err
			}

			if c.PPU.RenderDone || c.debug == DebugStepFrame {
				break
			}
		}
	}

	if c.debug != DebugDisabled {
		c.debug = DebugWait
	}

	return nil
}

func (c *Console) Draw(screen *ebiten.Image) {
	if c.PPU.RenderDone {
		img := c.PPU.Image()
		screen.WritePixels(img.Pix)
		c.PPU.RenderDone = false
	}
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
