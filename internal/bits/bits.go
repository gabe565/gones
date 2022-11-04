package bits

type Bits uint8

// Set sets bit(s)
func Set(b, flag Bits) Bits { return b | flag }

// Clear clears bit(s)
func Clear(b, flag Bits) Bits { return b &^ flag }

// Toggle toggles bit(s)
func Toggle(b, flag Bits) Bits { return b ^ flag }

// Has returns true if bit(s) set
func Has(b, flag Bits) bool { return b&flag != 0 }
