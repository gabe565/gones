package registers

// yyy NN YYYYY XXXXX
// ||| || ||||| +++++-- coarse X scroll
// ||| || +++++-------- coarse Y scroll
// ||| ++-------------- nametable select
// +++----------------- fine Y scroll

type Address struct {
	CoarseX    byte
	CoarseY    byte
	NametableX bool
	NametableY bool
	FineY      byte
}

func (r *Address) Get() uint16 {
	v := uint16(r.FineY&7)<<12 | uint16(r.CoarseY&31)<<5 | uint16(r.CoarseX&31)
	if r.NametableY {
		v |= 1 << 11
	}
	if r.NametableX {
		v |= 1 << 10
	}
	return v
}

func (r *Address) Set(data uint16) {
	r.CoarseX = byte(data & 31)
	r.CoarseY = byte(data >> 5 & 31)
	r.NametableX = data>>10&1 == 1
	r.NametableY = data>>11&1 == 1
	r.FineY = byte(data >> 12 & 7)
}

func (r *Address) WriteHi(data byte) {
	r.Set(uint16(data)<<8 | r.Get()&0xFF)
}

func (r *Address) WriteLo(data byte) {
	r.Set(r.Get()&0xFF00 | uint16(data))
}

func (r *Address) WriteScrollY(data byte) {
	r.CoarseY = data >> 3
	r.FineY = data & 7
}

func (r *Address) WriteScrollX(data byte) {
	r.CoarseX = data >> 3
}

func (r *Address) Increment(inc byte) {
	r.Set(r.Get() + uint16(inc))
}

func (r *Address) IncrementX() {
	if r.CoarseX < 31 {
		r.CoarseX++
	} else {
		r.CoarseX = 0
		r.NametableX = !r.NametableX
	}
}

func (r *Address) IncrementY() {
	if r.FineY < 7 {
		r.FineY++
	} else {
		r.FineY = 0

		switch r.CoarseY {
		case 29:
			r.CoarseY = 0
			r.NametableY = !r.NametableY
		case 31:
			r.CoarseY = 0
		default:
			r.CoarseY++
		}
	}
}

func (r *Address) LoadX(other Address) {
	r.NametableX = other.NametableX
	r.CoarseX = other.CoarseX
}

func (r *Address) LoadY(other Address) {
	r.NametableY = other.NametableY
	r.CoarseY = other.CoarseY
	r.FineY = other.FineY
}
