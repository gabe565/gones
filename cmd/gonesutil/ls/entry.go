package ls

import "github.com/gabe565/gones/internal/cartridge"

func newEntry(file string, cart *cartridge.Cartridge) *entry {
	return &entry{
		Path:    file,
		Name:    cart.Name(),
		Mapper:  cart.Mapper,
		Mirror:  cart.Mirror.String(),
		Battery: cart.Battery,
		Hash:    cart.Hash(),
	}
}

type entry struct {
	Path    string `json:"path" yaml:"path"`
	Name    string `json:"name" yaml:"name"`
	Mapper  byte   `json:"mapper" yaml:"mapper"`
	Mirror  string `json:"mirror" yaml:"mirror"`
	Battery bool   `json:"battery" yaml:"battery"`
	Hash    string `json:"hash" yaml:"hash"`
}
