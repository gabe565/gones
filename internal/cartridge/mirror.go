package cartridge

//go:generate go tool stringer -type Mirror

type Mirror byte

const (
	Horizontal Mirror = iota
	Vertical
	SingleLower
	SingleUpper
	FourScreen
)
