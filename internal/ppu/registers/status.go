package registers

// PPU Status bits
//
// 7 6 5 43210
// V S O .....
// ╷ ╷ ╷ ╷╷╷╷╷
// │ │ │ └┴┴┴┴╴PPU open bus. Returns stale PPU bus contents.
// │ │ └──────╴Sprite overflow. The intent was for this flag to be set
// │ │           whenever more than eight sprites appear on a scanline, but a
// │ │           hardware bug causes the actual behavior to be more complicated
// │ │           and generate false positives as well as false negatives; see
// │ │           PPU sprite evaluation. This flag is set during sprite
// │ │           evaluation and cleared at dot 1 (the second dot) of the
// │ │           pre-render line.
// │ └────────╴Sprite 0 Hit.  Set when a nonzero pixel of sprite 0 overlaps
// │             a nonzero background pixel; cleared at dot 1 of the pre-render
// │             line.  Used for raster timing.
// └──────────╴Vertical blank has started (0: not in vblank; 1: in vblank).
//               Set at dot 1 of line 241 (the line *after* the post-render
//               line); cleared after reading $2002 and at dot 1 of the
//               pre-render line.

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
