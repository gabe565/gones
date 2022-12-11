package registers

type Status struct {
	SpriteOverflow bool
	SpriteZeroHit  bool
	PrevVblank     bool
	Vblank         bool
}

const (
	StatusSpriteOverflow = 1 << (iota + 5)
	StatusSpriteZeroHit
	StatusVblank
)

func (s *Status) Set(data byte) {
	s.SpriteOverflow = data&StatusSpriteOverflow != 0
	s.SpriteZeroHit = data&StatusSpriteZeroHit != 0
	s.Vblank = data&StatusVblank != 0
}

func (s *Status) Get() byte {
	var v byte
	if s.SpriteOverflow {
		v |= StatusSpriteOverflow
	}
	if s.SpriteZeroHit {
		v |= StatusSpriteZeroHit
	}
	if s.Vblank {
		v |= StatusVblank
	}
	return v
}
