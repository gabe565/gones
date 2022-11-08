package registers

type Scroll struct {
	X     byte
	Y     byte
	Latch bool
}

func (s *Scroll) Write(data byte) {
	if !s.Latch {
		s.X = data
	} else {
		s.Y = data
	}
	s.Latch = !s.Latch
}

func (s *Scroll) ResetLatch() {
	s.Latch = false
}
