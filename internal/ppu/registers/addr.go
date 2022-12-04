package registers

// yyy NN YYYYY XXXXX
// ||| || ||||| +++++-- coarse X scroll
// ||| || +++++-------- coarse Y scroll
// ||| ++-------------- nametable select
// +++----------------- fine Y scroll

type AddrRegister struct {
	CoarseX    byte
	CoarseY    byte
	NametableX bool
	NametableY bool
	FineX      byte
	FineY      byte
	// Latch will write hi when false, lo when true
	Latch bool
}

func (r *AddrRegister) Get() uint16 {
	v := uint16(r.FineY&7)<<12 | uint16(r.CoarseY&31)<<5 | uint16(r.CoarseX&31)
	if r.NametableY {
		v |= 1 << 11
	}
	if r.NametableX {
		v |= 1 << 10
	}
	return v
}

func (r *AddrRegister) Set(data uint16) {
	r.CoarseX = byte(data & 31)
	r.CoarseY = byte(data >> 5 & 31)
	r.NametableX = data>>10&1 == 1
	r.NametableY = data>>11&1 == 1
	r.FineY = byte(data >> 12 & 7)
}

func (r *AddrRegister) Write(data byte) {
	v := r.Get()
	if r.Latch {
		// Lo
		r.Set(v&0xFF00 | uint16(data))
	} else {
		// Hi
		r.Set(uint16(data)<<8 | v&0xFF)
	}
	r.Latch = !r.Latch
}

func (r *AddrRegister) WriteScroll(data byte) {
	if r.Latch {
		r.CoarseY = data >> 3
		r.FineY = data & 7
	} else {
		r.CoarseX = data >> 3
		r.FineX = data & 7
	}
	r.Latch = !r.Latch
}

func (r *AddrRegister) Increment(inc byte) {
	r.Set(r.Get() + uint16(inc))
}

func (r *AddrRegister) ResetLatch() {
	r.Latch = false
}
