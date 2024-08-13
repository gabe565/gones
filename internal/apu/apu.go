package apu

import (
	"log/slog"

	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/interrupt"
	"github.com/gabe565/gones/internal/memory"
	"github.com/gabe565/gones/internal/util"
)

type CPU interface {
	memory.Read8
	interrupt.Stall
}

const (
	FrameCounterRate  = consts.CPUFrequency / 240.0
	DefaultSampleRate = consts.CPUFrequency / consts.AudioSampleRate * consts.FrameRateDifference
	BufferCap         = consts.AudioSampleRate / 5
)

//nolint:gochecknoglobals
var (
	lengthTable = [...]byte{
		10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
		12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
	}
	squareTable [31]float64
	tndTable    [203]float64
)

const (
	StatusPulse1 = 1 << iota
	StatusPulse2
	StatusTriangle
	StatusNoise
	StatusDMC
)

func init() { //nolint:all
	for i := range squareTable {
		squareTable[i] = 95.52 / (8128.0/float64(i) + 100)
	}
	for i := range tndTable {
		tndTable[i] = 163.67 / (24329.0/float64(i) + 100)
	}
}

func New(conf *config.Config) *APU {
	return &APU{
		Enabled:    true,
		SampleRate: DefaultSampleRate,
		conf:       &conf.Audio,

		Square: [2]Square{{Channel1: true}, {}},
		Noise:  Noise{ShiftRegister: 1},

		buf: make(chan float64, BufferCap),
	}
}

type APU struct {
	Enabled    bool    `msgpack:"-"`
	SampleRate float64 `msgpack:"-"`
	conf       *config.Audio

	Square   [2]Square
	Triangle Triangle
	Noise    Noise
	DMC      DMC

	Cycle       uint
	FramePeriod uint8
	FrameValue  byte

	IRQEnabled bool `msgpack:"alias:IrqEnabled"`
	IRQPending bool `msgpack:"alias:IrqPending"`

	buf chan float64
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
		if a.FramePeriod == 5 {
			a.stepEnvelope()
			a.stepSweep()
			a.stepLength()
		}
	default:
		slog.Error("Invalid APU write", "addr", util.EncodeHexAddr(addr))
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
		a.IRQPending = false
		return data
	default:
		return 0
	}
}

func (a *APU) Reset() {
	a.Square[0].SetEnabled(false)
	a.Square[1].SetEnabled(false)
	a.Triangle.SetEnabled(false)
	a.Noise.SetEnabled(false)
	a.DMC.SetEnabled(false)
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
	switch a.FramePeriod {
	case 4, 5:
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
}

func (a *APU) sendSample() {
	if len(a.buf) < BufferCap {
		a.buf <- a.output()
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

func (a *APU) output() float64 {
	var p1, p2 byte
	if a.conf.Channels.Square1 {
		p1 = a.Square[0].output()
	}
	if a.conf.Channels.Square2 {
		p2 = a.Square[1].output()
	}
	pulseOut := squareTable[p1+p2]

	var t, n, d byte
	if a.conf.Channels.Triangle {
		t = a.Triangle.output()
	}
	if a.conf.Channels.Noise {
		n = a.Noise.output()
	}
	if a.conf.Channels.PCM {
		d = a.DMC.output()
	}
	tndOut := tndTable[3*t+2*n+d]

	return pulseOut + tndOut
}

func (a *APU) Clear() {
	for len(a.buf) != 0 {
		<-a.buf
	}
}

func (a *APU) Read(p []byte) (int, error) {
	var i int
	for i = 0; i < len(p); i += 4 {
		p := p[i : i+4 : i+4]
		select {
		case sample := <-a.buf:
			out := int16(sample * 32767)
			lo := byte(out)
			p[0], p[2] = lo, lo
			hi := byte(out >> 8)
			p[1], p[3] = hi, hi
		default:
			p[0], p[1], p[2], p[3] = 0, 0, 0, 0
		}
	}
	return i, nil
}
