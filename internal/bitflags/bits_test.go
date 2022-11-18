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

func TestBits_Set(t *testing.T) {
	b := Flags(0)
	b.Set(1, true)
	assert.EqualValues(t, 1, b)
	b.Set(1, false)
	assert.EqualValues(t, 0, b)
}

func TestBits_Intersects(t *testing.T) {
	b := Flags(0b110)
	assert.Equal(t, true, b.Intersects(0b110))
	assert.Equal(t, true, b.Intersects(0b10))
	assert.Equal(t, false, b.Intersects(0b1))
	assert.Equal(t, true, b.Intersects(0b101))
}

func TestBits_Intersection(t *testing.T) {
	assert.EqualValues(t, 0b100, Flags(0b101).Intersection(0b110))
}

func TestBits_Union(t *testing.T) {
	assert.EqualValues(t, 0b110, Flags(0b010).Union(0b110))
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
