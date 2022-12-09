package cpu

// Processor Status bits
//
//	.----------------- Negative Flag
//	| .--------------- Overflow Flag
//	| |   .----------- Break Command
//	| |   | .--------- Decimal Mode (not used on NES)
//	| |   | | .------- Interrupt Disable
//	| |   | | | .----- Zero Flag
//	| |   | | | | .--- Carry Flag
//	N V _ B D I Z C
//	7 6 5 4 3 2 1 0
const (
	Carry = 1 << iota
	Zero
	InterruptDisable
	Decimal
	Break
	Break2
	Overflow
	Negative
)

var DefaultStatus = Status{InterruptDisable: true, Break2: true}

type Status struct {
	Carry            bool
	Zero             bool
	InterruptDisable bool
	Decimal          bool
	Break            bool
	Break2           bool
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
	if s.Break {
		v |= Break
	}
	if s.Break2 {
		v |= Break2
	}
	if s.Overflow {
		v |= Overflow
	}
	if s.Negative {
		v |= Negative
	}
	return v
}

func (s *Status) Set(data byte) {
	s.Carry = data&Carry != 0
	s.Zero = data&Zero != 0
	s.InterruptDisable = data&InterruptDisable != 0
	s.Decimal = data&Decimal != 0
	s.Break = data&Break != 0
	s.Break2 = data&Break2 != 0
	s.Overflow = data&Overflow != 0
	s.Negative = data&Negative != 0
}
