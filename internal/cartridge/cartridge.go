package cartridge

type Cartridge struct {
	Prg     []byte
	Chr     []byte
	Sram    []byte
	Mapper  byte
	Mirror  Mirror
	Battery byte
}

func New() *Cartridge {
	return &Cartridge{
		Sram: make([]byte, 0x2000),
	}
}

func FromBytes(b []byte) *Cartridge {
	cart := New()

	cart.Prg = make([]byte, 0x600, 0x8000)
	cart.Prg = append(cart.Prg, b...)
	cart.Prg = cart.Prg[:cap(cart.Prg)]
	cart.Prg[0xFFFD-0x8000] = 0x86

	cart.Chr = make([]byte, 0x2000)

	return cart
}
