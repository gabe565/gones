package registers

type Scroll struct {
	X, Y  byte
	Latch bool
}

func (s *Scroll) Write(data byte) {
	if s.Latch {
		s.Y = data
	} else {
		s.X = data
	}
	s.Latch = !s.Latch
}

func (s *Scroll) ResetLatch() {
	s.Latch = false
}
