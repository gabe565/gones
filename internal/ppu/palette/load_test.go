package palette

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed default.pal
var palFile []byte

func TestLoadPal(t *testing.T) {
	t.Parallel()
	UpdateEmphasized()

	var defaultPalette Palette
	copy(defaultPalette[:], Default[:])
	var emphasizeR Palette
	copy(emphasizeR[:], EmphasizeR[:])
	var emphasizeG Palette
	copy(emphasizeG[:], EmphasizeG[:])
	var emphasizeB Palette
	copy(emphasizeB[:], EmphasizeB[:])
	var emphasizeRG Palette
	copy(emphasizeRG[:], EmphasizeRG[:])
	var emphasizeRB Palette
	copy(emphasizeRB[:], EmphasizeRB[:])
	var emphasizeGB Palette
	copy(emphasizeGB[:], EmphasizeGB[:])
	var emphasizeRGB Palette
	copy(emphasizeRGB[:], EmphasizeRGB[:])

	Default = Palette{}
	EmphasizeR = Palette{}
	EmphasizeG = Palette{}
	EmphasizeB = Palette{}
	EmphasizeRG = Palette{}
	EmphasizeRB = Palette{}
	EmphasizeGB = Palette{}
	EmphasizeRGB = Palette{}

	if err := LoadPal(bytes.NewReader(palFile)); !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, defaultPalette, Default)
	assert.Equal(t, emphasizeR, EmphasizeR)
	assert.Equal(t, emphasizeG, EmphasizeG)
	assert.Equal(t, emphasizeB, EmphasizeB)
	assert.Equal(t, emphasizeRG, EmphasizeRG)
	assert.Equal(t, emphasizeRB, EmphasizeRB)
	assert.Equal(t, emphasizeGB, EmphasizeGB)
	assert.Equal(t, emphasizeRGB, EmphasizeRGB)
}
