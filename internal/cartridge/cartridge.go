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
