package cpu

//go:generate stringer -type InstructionName

type InstructionName uint8

const (
	BRK InstructionName = iota
	ORA
	SLO
	NOP
	ASL
	PHP
	ANC
	BPL
	CLC
	JSR
	AND
	RLA
	BIT
	ROL
	PLP
	BMI
	SEC
	RTI
	EOR
	SRE
	LSR
	PHA
	ALR
	JMP
	BVC
	CLI
	RTS
	ADC
	RRA
	ROR
	PLA
	ARR
	BVS
	SEI
	STA
	SAX
	STY
	STX
	DEY
	TXA
	XAA
	BCC
	AHX
	TYA
	TXS
	TAS
	SHY
	SHX
	LDY
	LDA
	LDX
	LAX
	TAY
	TAX
	LXA
	BCS
	CLV
	TSX
	LAS
	CPY
	CMP
	DCP
	DEC
	INY
	DEX
	AXS
	BNE
	CLD
	CPX
	SBC
	ISB
	INC
	INX
	BEQ
	SED
)
