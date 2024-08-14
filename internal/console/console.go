package console

import (
	"errors"
	"fmt"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/gabe565/gones/internal/ppu/palette"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const AutoSaveNum = 0

var ErrExit = errors.New("exit")

type UpdateAction uint8

const (
	ActionNone UpdateAction = iota
	ActionExit
	ActionSaveState
	ActionLoadState
)

type Console struct {
	config *config.Config

	CPU *cpu.CPU
	Bus *bus.Bus
	PPU *ppu.PPU
	APU *apu.APU

	Cartridge *cartridge.Cartridge
	Mapper    cartridge.Mapper

	audioCtx       *audio.Context
	player         *audio.Player
	actionOnUpdate UpdateAction
	enableTrace    bool
	debug          Debug

	undoSaveStates [][]byte
	undoLoadStates [][]byte

	autosave *time.Ticker
	rate     uint8

	willScreenshot bool
}

func New(conf *config.Config, cart *cartridge.Cartridge) (*Console, error) {
	console := Console{
		config:    conf,
		Cartridge: cart,
		rate:      1,

		undoSaveStates: make([][]byte, 0, conf.State.UndoStateCount),
		undoLoadStates: make([][]byte, 0, conf.State.UndoStateCount),
	}

	var err error
	console.Mapper, err = cartridge.NewMapper(cart)
	if err != nil {
		return &console, err
	}

	if err := console.LoadSRAM(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return &console, err
	}

	if err := palette.LoadPalFile(conf.UI.Palette); err != nil {
		return &console, err
	}

	console.PPU = ppu.New(console.Mapper)
	console.APU = apu.New(conf)
	console.Bus = bus.New(conf, console.Mapper, console.PPU, console.APU)
	console.CPU = cpu.New(console.Bus)

	console.PPU.SetCPU(console.CPU)
	console.APU.SetCPU(console.CPU)

	if conf.Audio.Enabled {
		console.audioCtx = audio.NewContext(consts.AudioSampleRate)
		console.player, err = console.audioCtx.NewPlayer(console.APU)
		if err != nil {
			return &console, err
		}
		console.player.SetBufferSize(time.Second / 20)
		console.player.SetVolume(conf.Audio.Volume)
		go func() {
			console.player.Play()
		}()
	} else {
		console.APU.Enabled = false
	}

	if conf.State.Resume {
		if err := console.LoadStateNum(AutoSaveNum); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return &console, err
			}
		}
	}

	console.SetTrace(conf.Debug.Trace)
	console.SetDebug(conf.Debug.Enabled)

	if duration := conf.State.AutosaveInterval; duration != 0 {
		console.autosave = time.NewTicker(time.Duration(duration))
	}

	return &console, nil
}

func (c *Console) Close() error {
	var errs []error
	if c.autosave != nil {
		c.autosave.Stop()
	}
	if c.config.State.Resume {
		errs = append(errs, c.SaveStateNum(AutoSaveNum, false))
	}
	errs = append(errs, c.SaveSRAM())
	return errors.Join(errs...)
}

func (c *Console) Step(render bool) {
	if c.enableTrace {
		//nolint:forbidigo
		fmt.Println(c.Trace())
	}

	var irq bool

	cycles := c.CPU.Step()
	if mapper, ok := c.Mapper.(cartridge.MapperOnCPUStep); ok {
		mapper.OnCPUStep(cycles)
	}

	for range cycles * 3 {
		c.PPU.Step(render)
	}

	for range cycles {
		irq = c.APU.Step() || irq
	}

	if mapper, ok := c.Mapper.(cartridge.MapperIRQ); ok {
		irq = mapper.IRQ() || irq
	}

	c.CPU.IRQPending = irq
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
	switch c.actionOnUpdate {
	case ActionNone:
	case ActionExit:
		return ErrExit
	case ActionSaveState:
		if err := c.SaveStateNum(1, true); err != nil {
			slog.Error("Failed to save state", "error", err)
		}
		c.actionOnUpdate = ActionNone
	case ActionLoadState:
		if err := c.LoadStateNum(1); err != nil {
			slog.Error("Failed to load state", "error", err)
		}
		c.actionOnUpdate = ActionNone
	}

	c.CheckInput()

	if c.debug == DebugWait {
		return nil
	}

	for i := range c.rate {
		if c.rate != 1 {
			c.PPU.RenderDone = false
		}
		for {
			c.Step(i == c.rate-1)

			if c.PPU.RenderDone || c.debug == DebugStepFrame {
				break
			}
		}
	}

	if c.debug != DebugDisabled {
		c.debug = DebugWait
	}

	if c.autosave != nil {
		select {
		case <-c.autosave.C:
			if err := c.SaveSRAM(); err != nil {
				slog.Error("Auto-save failed", "error", err)
			}
			if c.config.State.Resume {
				if err := c.SaveStateNum(AutoSaveNum, false); err != nil {
					slog.Error("State auto-save failed", "error", err)
				}
			}
		default:
		}
	}

	return nil
}

func (c *Console) Draw(screen *ebiten.Image) {
	if c.willScreenshot && runtime.GOOS != "js" {
		c.willScreenshot = false
		if err := c.writeScreenshot(screen); err != nil {
			slog.Error("Screenshot failed", "error", err)
		}
	}

	if c.PPU.RenderDone {
		img := c.PPU.Image()
		screen.WritePixels(img.Pix)
		c.PPU.RenderDone = false
	}
}

func (c *Console) SetUpdateAction(action UpdateAction) {
	c.actionOnUpdate = action
	if action == ActionExit {
		ebiten.SetRunnableOnUnfocused(true)
	}
}

func (c *Console) SetTrace(v bool) {
	c.enableTrace = v
}

func (c *Console) SetDebug(v bool) {
	if v {
		slog.Info("Enable step debug")
		c.debug = DebugWait
	} else {
		c.debug = DebugDisabled
	}
}

func (c *Console) writeScreenshot(screen *ebiten.Image) error {
	dir, err := config.GetScreenshotDir()
	if err != nil {
		return err
	}

	gameDir := filepath.Join(dir, c.Cartridge.Name())
	filename := filepath.Join(gameDir, time.Now().Format("2006-01-02_150405")+".png")

	if err := os.MkdirAll(gameDir, 0o777); err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err := png.Encode(f, screen); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	slog.Info("Saved screenshot", "path", filename)
	return nil
}
