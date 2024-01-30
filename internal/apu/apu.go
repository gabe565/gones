package apu

import (
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/memory"
	log "github.com/sirupsen/logrus"
)

type CPU interface {
	memory.Read8
	interrupts.Interruptible
	interrupts.Stallable
}

const (
	FrameCounterRate  = consts.CpuFrequency / 240.0
	DefaultSampleRate = consts.CpuFrequency / consts.AudioSampleRate * consts.FrameRateDifference
)

var lengthTable = [...]byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

const (
	StatusPulse1 = 1 << iota
	StatusPulse2
	StatusTriangle
	StatusNoise
	StatusDMC
)

func New() *APU {
	return &APU{
		Enabled:    true,
		SampleRate: DefaultSampleRate,

		Square: [2]Square{{Channel1: true}, {}},
		Noise:  Noise{ShiftRegister: 1},

		buf: make(chan float32, 8*consts.AudioSampleRate/20),
	}
}

type APU struct {
	Enabled    bool    `msgpack:"-"`
	SampleRate float64 `msgpack:"-"`

	Square   [2]Square
	Triangle Triangle
	Noise    Noise
	DMC      DMC

	Cycle       uint
	FramePeriod uint8
	FrameValue  byte

	IrqEnabled bool
	IrqPending bool

	buf chan float32
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
		a.DMC.IrqPending = false
	case addr == 0x4017:
		a.FramePeriod = 4 + data>>7&1
		a.IrqEnabled = data>>6&1 == 0
		if a.FramePeriod == 5 {
			a.stepEnvelope()
			a.stepSweep()
			a.stepLength()
		}
	default:
		log.Warnf("invalid APU write to $%04X", addr)
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
		a.IrqPending = false
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
	a.Cycle += 1
	cycle2 := float64(a.Cycle)

	a.stepTimer()

	f1 := uint32(cycle1 / FrameCounterRate)
	f2 := uint32(cycle2 / FrameCounterRate)
	if f1 != f2 {
		a.stepFrameCounter()
	}

	s1 := uint32(cycle1 / a.SampleRate)
	s2 := uint32(cycle2 / a.SampleRate)
	if s1 != s2 && a.Enabled {
		a.sendSample()
	}

	return a.IrqPending || a.DMC.IrqPending
}

func (a *APU) SetCpu(c CPU) {
	a.DMC.cpu = c
}

func (a *APU) stepFrameCounter() {
	switch a.FramePeriod {
	case 4, 5:
		a.FrameValue += 1
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
			if a.FramePeriod == 4 && a.IrqEnabled {
				a.IrqPending = true
			}
		}
	}
}

func (a *APU) sendSample() {
	select {
	case a.buf <- a.output():
	default:
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
	p1 := float32(a.Square[0].output())
	p2 := float32(a.Square[1].output())
	pulseOut := 95.88 / (8128/(p1+p2) + 100)

	t := float32(a.Triangle.output())
	n := float32(a.Noise.output())
	d := float32(a.DMC.output())
	tndOut := 159.79 / (1/(t/8227+n/12241+d/22638) + 100)

	return pulseOut + tndOut
}

func (a *APU) Clear() {
	for len(a.buf) != 0 {
		<-a.buf
	}
}

func (a *APU) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i += 4 {
		sample := <-a.buf
		p := p[i : i+4 : i+4]
		out := int16(sample * 32767)
		lo := byte(out)
		p[0], p[2] = lo, lo
		hi := byte(out >> 8)
		p[1], p[3] = hi, hi
	}
	return len(p), nil
}
