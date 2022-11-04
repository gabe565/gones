package bits

type Bits uint8

// Set sets bit(s) unconditionally
func (b *Bits) Set(flag Bits) { *b |= flag }

// Clear clears bit(s)
func (b *Bits) Clear(flag Bits) { *b &^= flag }

// Toggle toggles bit(s)
func (b *Bits) Toggle(flag Bits) { *b ^= flag }

// Has returns true if bit(s) set
func (b Bits) Has(flag Bits) bool { return b&flag != 0 }
