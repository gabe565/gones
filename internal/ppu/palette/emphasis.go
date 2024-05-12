package palette

import (
	"math"
)

//nolint:gochecknoglobals
var (
	EmphasizeR   = Palette{Emphasis: Red}
	EmphasizeG   = Palette{Emphasis: Green}
	EmphasizeB   = Palette{Emphasis: Blue}
	EmphasizeRG  = Palette{Emphasis: Red | Green}
	EmphasizeRB  = Palette{Emphasis: Red | Blue}
	EmphasizeGB  = Palette{Emphasis: Green | Blue}
	EmphasizeRGB = Palette{Emphasis: Red | Green | Blue}
)

type Emphasis uint8

const (
	Red Emphasis = 1 << iota
	Green
	Blue
)

const Attenuate = 0.746

func UpdateEmphasized() {
	palettes := []*Palette{
		&EmphasizeR,
		&EmphasizeG,
		&EmphasizeB,
		&EmphasizeRG,
		&EmphasizeRB,
		&EmphasizeGB,
		&EmphasizeRGB,
	}

	for _, palette := range palettes {
		for i, c := range Default.RGBA {
			// Don't attenuate $xE or $xF (black)
			if i&0xE != 0xE {
				if palette.Emphasis&Red != 0 {
					c.G = uint8(math.Round(float64(c.G) * Attenuate))
					c.B = uint8(math.Round(float64(c.B) * Attenuate))
				}
				if palette.Emphasis&Green != 0 {
					c.R = uint8(math.Round(float64(c.R) * Attenuate))
					c.B = uint8(math.Round(float64(c.B) * Attenuate))
				}
				if palette.Emphasis&Blue != 0 {
					c.R = uint8(math.Round(float64(c.R) * Attenuate))
					c.G = uint8(math.Round(float64(c.G) * Attenuate))
				}
			}
			palette.RGBA[i] = c
		}
	}
}
