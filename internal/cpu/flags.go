package cpu

import "github.com/gabe565/gones/internal/bits"

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
	Carry bits.Bits = 1 << iota
	Zero
	InterruptDisable
	DecimalMode
	Break
	Break2
	Overflow
	Negative
)
