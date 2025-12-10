package cpu

// Processor Status bits
//
// 7 6 5 4 3 2 1 0
// N V _ B D I Z C
// ╷ ╷   ╷ ╷ ╷ ╷ ╷
// │ │   │ │ │ │ └╴Carry Flag
// │ │   │ │ │ └──╴Zero Flag
// │ │   │ | └────╴Interrupt Disable
// │ │   │ └──────╴Decimal Mode (not used on NES)
// │ │   └────────╴Break Flag
// │ │
// │ └────────────╴Overflow Flag
// └──────────────╴Negative Flag

const (
	Carry = 1 << iota
	Zero
	InterruptDisable
	Decimal
	Break
	Unused
	Overflow
	Negative
)

type Status struct {
	Carry            bool
	Zero             bool
	InterruptDisable bool
	Decimal          bool
	Break            bool
	Overflow         bool
	Negative         bool
}

func (s *Status) Get() byte {
	var v byte
	if s.Carry {
		v |= Carry
	}
	if s.Zero {
		v |= Zero
	}
	if s.InterruptDisable {
		v |= InterruptDisable
	}
	if s.Decimal {
		v |= Decimal
	}
	if s.Overflow {
		v |= Overflow
	}
	if s.Negative {
		v |= Negative
	}
	return v | Unused
}

func (s *Status) Set(data byte) {
	s.Carry = data&Carry != 0
	s.Zero = data&Zero != 0
	s.InterruptDisable = data&InterruptDisable != 0
	s.Decimal = data&Decimal != 0
	s.Overflow = data&Overflow != 0
	s.Negative = data&Negative != 0
}
