package memory

type Read8 interface {
	ReadMem(uint16) byte
}

type Write8 interface {
	WriteMem(uint16, byte)
}

type ReadWrite8 interface {
	Read8
	Write8
}

type Read16 interface {
	ReadMem16(uint16) uint16
}

type Write16 interface {
	WriteMem16(uint16, uint16)
}

type ReadWrite16 interface {
	Read16
	Write16
}

type ReadWrite interface {
	ReadWrite8
	ReadWrite16
}
