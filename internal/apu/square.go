package apu

var duties = [...][8]byte{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

type Square struct {
	Enabled  bool
	Channel1 bool

	DutyMode  byte
	DutyValue byte

	EnvelopeEnabled bool
	EnvelopePeriod  byte
	EnvelopeLoop    bool
	EnvelopeStart   bool
	EnvelopeVol     byte
	EnvelopeValue   byte

	Volume byte

	SweepEnabled bool
	SweepPeriod  byte
	SweepNegate  bool
	SweepShift   byte
	SweepReload  bool
	SweepValue   byte

	LengthEnabled bool
	LengthValue   byte

	TimerPeriod uint16
	TimerValue  uint16
}

func (p *Square) Write(addr uint16, data byte) {
	switch addr {
	case 0x4000, 0x4004:
		p.DutyMode = data >> 6 & 3
		p.LengthEnabled = data>>5&1 == 0
		p.EnvelopeLoop = data>>5&1 == 1
		p.EnvelopeEnabled = data>>4&1 == 0
		p.Volume = data & 0xF
		p.EnvelopePeriod = data & 0xF
	case 0x4001, 0x4005:
		p.SweepEnabled = data>>7&1 == 1
		p.SweepPeriod = data >> 4 & 7
		p.SweepNegate = data>>3&1 == 1
		p.SweepShift = data & 0x7
		p.SweepReload = true
	case 0x4002, 0x4006:
		p.TimerPeriod = p.TimerPeriod&0x700 | uint16(data)
	case 0x4003, 0x4007:
		if p.Enabled {
			p.LengthValue = lengths[data>>3&0x1F]
		}
		p.TimerPeriod = uint16(data)&0x7<<8 | p.TimerPeriod&0xFF
		p.EnvelopeStart = true
		p.DutyValue = 0
	}
}

func (p *Square) SetEnabled(v bool) {
	p.Enabled = v
	if !v {
		p.LengthValue = 0
	}
}

func (p *Square) stepTimer() {
	if p.TimerValue == 0 {
		p.TimerValue = p.TimerPeriod
		p.DutyValue += 1
		p.DutyValue %= 8
	} else {
		p.TimerValue -= 1
	}
}

func (p *Square) stepEnvelope() {
	if p.EnvelopeStart {
		p.EnvelopeVol = 15
		p.EnvelopeValue = p.EnvelopePeriod
		p.EnvelopeStart = false
	} else if p.EnvelopeValue > 0 {
		p.EnvelopeValue -= 1
	} else {
		if p.EnvelopeVol > 0 {
			p.EnvelopeVol -= 1
		} else if p.EnvelopeLoop {
			p.EnvelopeVol = 15
		}
		p.EnvelopeValue = p.EnvelopePeriod
	}
}

func (p *Square) stepSweep() {
	if p.SweepReload {
		if p.SweepEnabled && p.SweepValue == 0 {
			p.sweep()
		}
		p.SweepValue = p.SweepPeriod
		p.SweepReload = false
	} else if p.SweepValue > 0 {
		p.SweepValue -= 1
	} else {
		if p.SweepEnabled {
			p.sweep()
		}
		p.SweepValue = p.SweepPeriod
	}
}

func (p *Square) sweep() {
	delta := p.TimerPeriod >> p.SweepShift
	if p.SweepNegate {
		delta = -delta
		if p.Channel1 {
			delta -= 1
		}
	}
	p.TimerPeriod += delta
}

func (p *Square) stepLength() {
	if p.LengthEnabled && p.LengthValue > 0 {
		p.LengthValue -= 1
	}
}

func (p *Square) output() byte {
	switch {
	case !p.Enabled:
		return 0
	case p.LengthValue == 0:
		return 0
	case duties[p.DutyMode][p.DutyValue] == 0:
		return 0
	case p.TimerPeriod < 8, p.TimerPeriod > 0x7FF:
		return 0
	case p.EnvelopeEnabled:
		return p.EnvelopeVol
	default:
		return p.Volume
	}
}
