package apu

var noisePeriodTable = [...]uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}

type Noise struct {
	Enabled bool

	EnvelopeEnabled bool
	EnvelopeLoop    bool
	EnvelopeStart   bool
	EnvelopeVol     byte
	EnvelopePeriod  byte
	EnvelopeValue   byte

	Volume byte

	LoopNoise     bool
	ShiftRegister uint16

	TimerPeriod uint16
	TimerValue  uint16

	LengthEnabled bool
	LengthValue   byte
}

func (n *Noise) Write(addr uint16, data byte) {
	switch addr {
	case 0x400C:
		n.EnvelopeLoop = data>>5&1 == 1
		n.LengthEnabled = data>>5&1 == 0
		n.EnvelopeEnabled = data>>4&1 == 0
		n.EnvelopePeriod = data & 0xF
		n.Volume = data & 0xF
	case 0x400D:
		//
	case 0x400E:
		n.LoopNoise = data>>7&1 == 1
		n.TimerPeriod = noisePeriodTable[data&0xF]
	case 0x400F:
		if n.Enabled {
			n.LengthValue = lengthTable[data>>3&0x1F]
		}
		n.EnvelopeStart = true
	}
}

func (n *Noise) SetEnabled(v bool) {
	n.Enabled = v
	if !v {
		n.LengthValue = 0
	}
}

func (n *Noise) stepTimer() {
	if n.TimerValue == 0 {
		n.TimerValue = n.TimerPeriod
		feedback := n.ShiftRegister & 1
		if n.LoopNoise {
			feedback ^= n.ShiftRegister >> 6 & 1
		} else {
			feedback ^= n.ShiftRegister >> 1 & 1
		}
		n.ShiftRegister >>= 1
		n.ShiftRegister |= feedback << 14
	} else {
		n.TimerValue -= 1
	}
}

func (n *Noise) stepEnvelope() {
	if n.EnvelopeStart {
		n.EnvelopeVol = 15
		n.EnvelopeValue = n.EnvelopePeriod
		n.EnvelopeStart = false
	} else if n.EnvelopeValue > 0 {
		n.EnvelopeValue -= 1
	} else {
		if n.EnvelopeVol > 0 {
			n.EnvelopeVol -= 1
		} else if n.EnvelopeLoop {
			n.EnvelopeVol = 15
		}
		n.EnvelopeValue = n.EnvelopePeriod
	}
}

func (n *Noise) stepLength() {
	if n.LengthEnabled && n.LengthValue > 0 {
		n.LengthValue -= 1
	}
}

func (n *Noise) output() byte {
	switch {
	case !n.Enabled:
		return 0
	case n.LengthValue == 0:
		return 0
	case n.ShiftRegister&1 == 1:
		return 0
	case n.EnvelopeEnabled:
		return n.EnvelopeVol
	default:
		return n.Volume
	}
}
