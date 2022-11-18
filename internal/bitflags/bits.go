package bitflags

type Flags byte

// Insert inserts the specified flags in-place.
func (f *Flags) Insert(other Flags) { *f |= other }

// Remove removes the specified flags in-place.
func (f *Flags) Remove(other Flags) { *f &^= other }

// Set inserts or removes the specified flags depending on the passed value.
func (f *Flags) Set(other Flags, value bool) {
	if value {
		f.Insert(other)
	} else {
		f.Remove(other)
	}
}

// Contains returns `true` if all of the flags in the parameter are contained in the current Flag
func (f Flags) Contains(other Flags) bool { return f&other == other }

// Intersects returns `true` if there are flags common to both the parameter and the current Flag.
func (f Flags) Intersects(other Flags) bool { return f&other != 0 }

// Intersection returns a new set of flags, containing only the flags present in both Flag and `other`.
func (f Flags) Intersection(other Flags) Flags { return f & other }

// Union returns a new set of flags, containing any flags present in either Flag or `other`.
func (f Flags) Union(other Flags) Flags { return f | other }

// Toggle the specified flags will be inserted if not present, and removed if they are.
func (f *Flags) Toggle(other Flags) { *f ^= other }
