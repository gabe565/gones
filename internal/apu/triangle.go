package apu

type Triangle struct {
	Enabled bool

	LengthEnabled bool
	LengthValue   byte

	CounterPeriod byte
	CounterValue  byte
	CounterReload bool

	DutyValue byte

	TimerPeriod uint16
	TimerValue  uint16
}

func (t *Triangle) Write(addr uint16, data byte) {
	switch addr {
	case 0x4008:
		t.LengthEnabled = data>>7&1 == 0
		t.CounterPeriod = data & 0x7F
	case 0x4009:
		//
	case 0x400A:
		t.TimerPeriod = t.TimerPeriod&0xFF00 | uint16(data)
	case 0x400B:
		t.LengthValue = lengths[data>>3&0x1F]
		t.TimerPeriod = uint16(data)&0x7<<8 | t.TimerPeriod&0xFF
		t.TimerValue = t.TimerPeriod
		t.CounterReload = true
	}
}

func (t *Triangle) SetEnabled(v bool) {
	t.Enabled = v
	if !v {
		t.LengthValue = 0
	}
}

func (t *Triangle) stepTimer() {
	if t.TimerValue == 0 {
		t.TimerValue = t.TimerPeriod
		if t.LengthValue > 0 && t.CounterValue > 0 {
			t.DutyValue += 1
			t.DutyValue %= 32
		}
	} else {
		t.TimerValue -= 1
	}
}

func (t *Triangle) stepLength() {
	if t.LengthEnabled && t.LengthValue > 0 {
		t.LengthValue -= 1
	}
}

func (t *Triangle) stepCounter() {
	if t.CounterReload {
		t.CounterValue = t.CounterPeriod
	} else if t.CounterValue > 0 {
		t.CounterValue -= 1
	}
	if t.LengthEnabled {
		t.CounterReload = false
	}
}

func (t *Triangle) output() byte {
	switch {
	case !t.Enabled:
		return 0
	case t.TimerPeriod < 3:
		return 0
	case t.LengthValue == 0:
		return 0
	case t.CounterValue == 0:
		return 0
	default:
		if t.DutyValue < 16 {
			return 15 - t.DutyValue
		} else {
			return t.DutyValue - 16
		}
	}
}
