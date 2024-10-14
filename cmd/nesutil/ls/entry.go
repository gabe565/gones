package ls

import "gabe565.com/gones/internal/cartridge"

func newEntry(file string, cart *cartridge.Cartridge) *entry {
	return &entry{
		Path:    file,
		Name:    cart.Name(),
		Mapper:  cart.Header.Mapper(),
		Mirror:  cart.Mirror.String(),
		Battery: cart.Battery,
		Hash:    cart.Hash(),
	}
}

type entry struct {
	Path    string `json:"path" yaml:"path"`
	Name    string `json:"name" yaml:"name"`
	Mapper  uint8  `json:"mapper" yaml:"mapper"`
	Mirror  string `json:"mirror" yaml:"mirror"`
	Battery bool   `json:"battery" yaml:"battery"`
	Hash    string `json:"hash" yaml:"hash"`
}
