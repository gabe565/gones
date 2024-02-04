package ls

import "github.com/gabe565/gones/internal/cartridge"

func newEntry(file string, cart *cartridge.Cartridge) entry {
	return entry{
		Path:    file,
		Name:    cart.Name(),
		Mapper:  cart.Mapper,
		Mirror:  cart.Mirror.String(),
		Battery: cart.Battery,
	}
}

type entry struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Mapper  byte   `json:"mapper"`
	Mirror  string `json:"mirror"`
	Battery bool   `json:"battery"`
}
