package apu

import (
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/interrupts"
	log "github.com/sirupsen/logrus"
)

const FrameCounterRate = consts.CpuFrequency / 240.0

var lengths = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

var squareTable [31]float32
var tndTable [203]float32

const (
	StatusPulse1 = 1 << iota
	StatusPulse2
	StatusTriangle
	StatusNoise
	StatusDMC
)

func init() {
	for i := range squareTable {
		squareTable[i] = 95.52 / (8128.0/float32(i) + 100)
	}
	for i := range tndTable {
		tndTable[i] = 163.67 / (24329.0/float32(i) + 100)
	}
}

func New() *APU {
	return &APU{
		Enabled:    true,
		Volume:     1,
		SampleRate: consts.CpuFrequency / float64(consts.AudioSampleRate),

		Square: [2]Square{{Channel: 1}, {Channel: 2}},
		Noise:  Noise{ShiftRegister: 1},

		buf: make(chan byte, 10*4*consts.AudioSampleRate/60),
	}
}

type APU struct {
	Enabled    bool
	SampleRate float64
	Volume     float32

	Square   [2]Square
	Triangle Triangle
	Noise    Noise
	DMC      DMC

	Cycle       uint
	FramePeriod uint8
	FrameValue  byte

	InterruptInhibit bool

	buf chan byte
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
	case addr == 0x4017:
		a.FramePeriod = 4 + data>>7&1
		a.InterruptInhibit = data>>6&1 == 1
		if a.FramePeriod == 5 {
			a.stepEnvelope()
			a.stepSweep()
			a.stepLength()
		}
	default:
		log.Fatalf("invalid APU write to $%04X", addr)
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
		return data
	default:
		return 0
	}
}

func (a *APU) Reset() {
	a.clearBuf()
	a.Square[0] = Square{Channel: 1}
	a.Square[1] = Square{Channel: 2}
	a.Triangle = Triangle{}
	a.Noise = Noise{ShiftRegister: 1}
	a.DMC = DMC{cpu: a.DMC.cpu}
}

func (a *APU) Step() *interrupts.Interrupt {
	var interrupt *interrupts.Interrupt

	cycle1 := a.Cycle
	a.Cycle += 1
	cycle2 := a.Cycle

	a.stepTimer()

	f1 := int(float32(cycle1) / FrameCounterRate)
	f2 := int(float32(cycle2) / FrameCounterRate)
	if f1 != f2 {
		interrupt = a.stepFrameCounter()
	}

	s1 := int(float64(cycle1) / a.SampleRate)
	s2 := int(float64(cycle2) / a.SampleRate)
	if s1 != s2 && a.Enabled {
		a.sendSample()
	}

	return interrupt
}

func (a *APU) SetCpu(c CPU) {
	a.DMC.cpu = c
}

func (a *APU) stepFrameCounter() *interrupts.Interrupt {
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
			if a.FramePeriod == 4 && !a.InterruptInhibit {
				return &interrupts.IRQ
			}
		}
	}
	return nil
}

func (a *APU) sendSample() {
	for _, b := range a.output() {
		select {
		case a.buf <- b:
		default:
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

func (a *APU) output() []byte {
	p1 := a.Square[0].output()
	p2 := a.Square[1].output()
	pulseOut := squareTable[p1+p2]

	t := a.Triangle.output()
	n := a.Noise.output()
	d := a.DMC.output()
	tndOut := tndTable[3*t+2*n+d]

	out := int16((pulseOut + tndOut) * a.Volume * 32767)
	return []byte{
		byte(out),
		byte(out >> 8),
		byte(out),
		byte(out >> 8),
	}
}

func (a *APU) clearBuf() {
	for {
		select {
		case <-a.buf:
		default:
			return
		}
	}
}

func (a *APU) Read(p []byte) (int, error) {
	var n int

loop:
	for i := range p {
		select {
		case p[i] = <-a.buf:
			n += 1
		default:
			break loop
		}
	}

	if n == 0 {
		for i := range p {
			p[i] = 0
		}
		n = len(p)
	}

	return n, nil
}
