package apu

import (
	"log/slog"
	"math"

	"gabe565.com/gones/internal/config"
	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/interrupt"
	"gabe565.com/gones/internal/log"
	"gabe565.com/gones/internal/memory"
)

type CPU interface {
	memory.Read8
	interrupt.Stall
}

const (
	FrameCounterRate  = float64(consts.CPUFrequency) / 240.0
	DefaultSampleRate = float64(consts.CPUFrequency) / float64(consts.AudioSampleRate) * consts.FrameRateDifference
	BufferCap         = consts.AudioSampleRate / 5 * 8
)

//nolint:gochecknoglobals
var (
	lengthTable = [...]byte{
		10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
		12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
	}
	squareTable [31]float32
	tndTable    [203]float32
)

const (
	StatusPulse1 = 1 << iota
	StatusPulse2
	StatusTriangle
	StatusNoise
	StatusDMC
	_
	StatusFrameInterrupt
	StatusDMCInterrupt
)

func init() { //nolint:all
	for i := range squareTable {
		squareTable[i] = float32(95.52 / (8128.0/float64(i) + 100))
	}
	for i := range tndTable {
		tndTable[i] = float32(163.67 / (24329.0/float64(i) + 100))
	}
}

func New(conf *config.Config) *APU {
	a := &APU{
		Enabled:    true,
		SampleRate: DefaultSampleRate,
		conf:       &conf.Audio,
		buf:        newRingBuffer(BufferCap),

		Square: [2]Square{{Channel1: true}, {}},
		Noise:  Noise{ShiftRegister: 1},

		FramePeriod: 4,
	}
	return a
}

type APU struct {
	Enabled    bool    `msgpack:"-"`
	SampleRate float64 `msgpack:"-"`
	conf       *config.Audio
	buf        *ringBuffer

	Square   [2]Square
	Triangle Triangle
	Noise    Noise
	DMC      DMC

	Cycle       uint
	FramePeriod uint8
	FrameValue  byte

	IRQEnabled bool `msgpack:"alias:IrqEnabled"`
	IRQPending bool `msgpack:"alias:IrqPending"`
}

func (a *APU) WriteMem(addr uint16, data byte) {
	switch {
	case 0x4000 <= addr && addr <= 0x4003:
		a.Square[0].Write(addr, data)
	case 0x4004 <= addr && addr <= 0x4007:
		a.Square[1].Write(addr, data)
	case 0x4008 <= addr && addr <= 0x400B:
		a.Triangle.Write(addr, data)
	case 0x400C <= addr && addr <= 0x400F:
		a.Noise.Write(addr, data)
	case 0x4010 <= addr && addr <= 0x4013:
		a.DMC.Write(addr, data)
	case addr == 0x4015:
		a.Square[0].SetEnabled(data&StatusPulse1 != 0)
		a.Square[1].SetEnabled(data&StatusPulse2 != 0)
		a.Triangle.SetEnabled(data&StatusTriangle != 0)
		a.Noise.SetEnabled(data&StatusNoise != 0)
		a.DMC.SetEnabled(data&StatusDMC != 0)
		a.DMC.IRQPending = false
	case addr == 0x4017:
		a.FramePeriod = 4 + data>>7&1
		a.IRQEnabled = data>>6&1 == 0
		if !a.IRQEnabled {
			a.IRQPending = false
		}
		if a.FramePeriod == 5 {
			a.stepEnvelope()
			a.stepSweep()
			a.stepLength()
		}
	default:
		slog.Error("Invalid APU write", "addr", log.HexAddr(addr))
	}
}

func (a *APU) ReadMem(addr uint16) byte {
	switch addr {
	case 0x4015:
		var data byte
		if a.Square[0].LengthValue > 0 {
			data |= StatusPulse1
		}
		if a.Square[1].LengthValue > 0 {
			data |= StatusPulse2
		}
		if a.Triangle.LengthValue > 0 {
			data |= StatusTriangle
		}
		if a.Noise.LengthValue > 0 {
			data |= StatusNoise
		}
		if a.DMC.CurrLen > 0 {
			data |= StatusDMC
		}
		if a.IRQPending {
			data |= StatusFrameInterrupt
		}
		if a.DMC.IRQPending {
			data |= StatusDMCInterrupt
		}
		a.IRQPending = false
		return data
	default:
		return 0
	}
}

func (a *APU) Reset() {
	a.IRQPending = false
	a.WriteMem(0x4015, 0)
}

func (a *APU) Step() bool {
	cycle1 := float64(a.Cycle)
	a.Cycle++
	cycle2 := float64(a.Cycle)

	a.stepTimer()

	f1 := uint32(cycle1 / FrameCounterRate)
	f2 := uint32(cycle2 / FrameCounterRate)
	if f1 != f2 {
		a.stepFrameCounter()
	}

	if a.Enabled {
		s1 := uint32(cycle1 / a.SampleRate)
		s2 := uint32(cycle2 / a.SampleRate)
		if s1 != s2 {
			a.sendSample()
		}
	}

	return a.IRQPending || a.DMC.IRQPending
}

func (a *APU) SetCPU(c CPU) {
	a.DMC.cpu = c
}

func (a *APU) stepFrameCounter() {
	a.FrameValue++
	a.FrameValue %= a.FramePeriod
	switch a.FrameValue {
	case 0, 2:
		a.stepEnvelope()
	case 1:
		a.stepEnvelope()
		a.stepSweep()
		a.stepLength()
	case 3:
		a.stepEnvelope()
		a.stepSweep()
		a.stepLength()
		if a.FramePeriod == 4 && a.IRQEnabled {
			a.IRQPending = true
		}
	}
}

func (a *APU) stepTimer() {
	if a.Cycle%2 == 0 {
		a.Square[0].stepTimer()
		a.Square[1].stepTimer()
		a.Noise.stepTimer()
		a.DMC.stepTimer()
	}
	a.Triangle.stepTimer()
}

func (a *APU) stepEnvelope() {
	a.Square[0].stepEnvelope()
	a.Square[1].stepEnvelope()
	a.Triangle.stepCounter()
	a.Noise.stepEnvelope()
}

func (a *APU) stepSweep() {
	a.Square[0].stepSweep()
	a.Square[1].stepSweep()
}

func (a *APU) stepLength() {
	a.Square[0].stepLength()
	a.Square[1].stepLength()
	a.Triangle.stepLength()
	a.Noise.stepLength()
}

func (a *APU) output() float32 {
	var square byte
	if a.conf.Channels.Square1 {
		square += a.Square[0].output()
	}
	if a.conf.Channels.Square2 {
		square += a.Square[1].output()
	}

	var tnd byte
	if a.conf.Channels.Triangle {
		tnd += 3 * a.Triangle.output()
	}
	if a.conf.Channels.Noise {
		tnd += 2 * a.Noise.output()
	}
	if a.conf.Channels.PCM {
		tnd += a.DMC.output()
	}

	return squareTable[square] + tndTable[tnd]
}

func (a *APU) sendSample() {
	result := a.output()
	b := math.Float32bits(result)
	a.buf.Write([]byte{
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
	})
}

func (a *APU) Clear() {
	a.buf.Reset()
}

func (a *APU) Read(p []byte) (int, error) {
	n := a.buf.Read(p)
	if n == 0 {
		clear(p)
		return len(p), nil
	}
	return n, nil
}
