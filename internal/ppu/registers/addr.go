package registers

type AddrRegister struct {
	Hi, Lo byte
	LoPtr  bool
}

func (r *AddrRegister) Get() uint16 {
	return uint16(r.Hi)<<8 | uint16(r.Lo)
}

func (r *AddrRegister) Set(data uint16) {
	r.Hi = uint8(data >> 8)
	r.Lo = uint8(data & 0xFF)
}

func (r *AddrRegister) Update(data byte) {
	if r.LoPtr {
		r.Lo = data
	} else {
		r.Hi = data
	}

	if v := r.Get(); v > 0x3FFF {
		r.Set(v & 0x3FFF)
	}
	r.LoPtr = !r.LoPtr
}

func (r *AddrRegister) Increment(inc byte) {
	lo := r.Lo
	r.Lo += inc
	if lo > r.Lo {
		r.Hi += 1
	}

	if v := r.Get(); v > 0x3FFF {
		r.Set(v & 0x3FFF)
	}
}

func (r *AddrRegister) ResetLatch() {
	r.LoPtr = false
}
