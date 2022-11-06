package bitflags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBits_Insert(t *testing.T) {
	b := Flags(0)
	b.Insert(0b10)
	assert.EqualValues(t, 0b10, b)
	b.Insert(0b1100)
	assert.EqualValues(t, 0b1110, b)
}

func TestBits_Remove(t *testing.T) {
	b := Flags(0b1110)
	b.Remove(0b10)
	assert.EqualValues(t, 0b1100, b)
	b.Remove(0b1100)
	assert.EqualValues(t, 0, b)

}

func TestBits_Toggle(t *testing.T) {
	b := Flags(0b10)
	b.Toggle(0b1)
	assert.EqualValues(t, 0b11, b)
	b.Toggle(0b10)
	assert.EqualValues(t, 0b1, b)
	b.Toggle(0b110)
	assert.EqualValues(t, 0b111, b)
}

func TestBits_Has(t *testing.T) {
	b := Flags(0b110)
	assert.Equal(t, true, b.Has(0b110))
	assert.Equal(t, true, b.Has(0b10))
	assert.Equal(t, false, b.Has(0b1))
}

func TestBits_Set(t *testing.T) {
	b := Flags(0)
	b.Set(1, true)
	assert.EqualValues(t, 1, b)
	b.Set(1, false)
	assert.EqualValues(t, 0, b)
}
