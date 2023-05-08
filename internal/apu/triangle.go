package apu

var triangleOutputTable = [...]byte{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
}

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
		if t.Enabled {
			t.LengthValue = lengthTable[data>>3&0x1F]
		}
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
		if t.LengthValue > 0 && t.CounterValue > 0 && t.TimerPeriod != 0 {
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
	return triangleOutputTable[t.DutyValue]
}
