package palette

import (
	"math"
)

var (
	EmphasizeR   Palette
	EmphasizeG   Palette
	EmphasizeB   Palette
	EmphasizeRG  Palette
	EmphasizeRB  Palette
	EmphasizeGB  Palette
	EmphasizeRGB Palette
)

const (
	Red   = 'R'
	Green = 'G'
	Blue  = 'B'

	Attenuate = 0.746
)

func init() {
	UpdateEmphasized()
}

func UpdateEmphasized() {
	type EmphasizedPalette struct {
		Palette *Palette
		Colors  []byte
	}

	palettes := []EmphasizedPalette{
		{&EmphasizeR, []byte{Red}},
		{&EmphasizeG, []byte{Green}},
		{&EmphasizeB, []byte{Blue}},
		{&EmphasizeRG, []byte{Red, Green}},
		{&EmphasizeRB, []byte{Red, Blue}},
		{&EmphasizeGB, []byte{Green, Blue}},
		{&EmphasizeRGB, []byte{Red, Green, Blue}},
	}

	for _, palette := range palettes {
		for i, c := range Default {
			// Don't attenuate $xE or $xF (black)
			if i&0xE != 0xE {
				for _, emphasis := range palette.Colors {
					switch emphasis {
					case Red:
						c.G = uint8(math.Round(float64(c.G) * Attenuate))
						c.B = uint8(math.Round(float64(c.B) * Attenuate))
					case Green:
						c.R = uint8(math.Round(float64(c.R) * Attenuate))
						c.B = uint8(math.Round(float64(c.B) * Attenuate))
					case Blue:
						c.R = uint8(math.Round(float64(c.R) * Attenuate))
						c.G = uint8(math.Round(float64(c.G) * Attenuate))
					}
				}
			}
			palette.Palette[i] = c
		}
	}
}
