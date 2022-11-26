package registers

type AddrRegister struct {
	Value uint16
	// Latch will write hi when false, lo when true
	Latch bool
}

func (r *AddrRegister) Get() uint16 {
	return r.Value
}

func (r *AddrRegister) Set(data uint16) {
	r.Value = data & 0x3FFF
}

func (r *AddrRegister) Write(data byte) {
	if r.Latch {
		// Lo
		r.Value &= 0xFF00
		r.Value |= uint16(data)
	} else {
		// Hi
		r.Value &= 0x00FF
		r.Value |= uint16(data) << 8
	}
	r.Value &= 0x3FFF
	r.Latch = !r.Latch
}

func (r *AddrRegister) Increment(inc byte) {
	r.Value += uint16(inc)
	r.Value &= 0x3FFF
}

func (r *AddrRegister) ResetLatch() {
	r.Latch = false
}
