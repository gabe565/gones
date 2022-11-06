package bitflags

type Flags byte

// Insert inserts flags in-place
func (b *Flags) Insert(flag Flags) { *b |= flag }

// Remove removes flags in-place
func (b *Flags) Remove(flag Flags) { *b &^= flag }

// Toggle specified flags will be inserted if not present, and removed if they are
func (b *Flags) Toggle(flag Flags) { *b ^= flag }

// Has returns true if flags are set
func (b Flags) Has(flag Flags) bool { return b&flag != 0 }

// Set inserts or removes the specified flags depending on the passed value
func (b *Flags) Set(flag Flags, condition bool) {
	if condition {
		b.Insert(flag)
	} else {
		b.Remove(flag)
	}
}
