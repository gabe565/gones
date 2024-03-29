package memory

type Read8 interface {
	ReadMem(addr uint16) byte
}

type ReadSafe interface {
	ReadMemSafe(addr uint16) byte
}

type Write8 interface {
	WriteMem(addr uint16, data byte)
}

type ReadWrite8 interface {
	Read8
	Write8
}

type Read16 interface {
	ReadMem16(addr uint16) uint16
}

type Write16 interface {
	WriteMem16(addr uint16, data uint16)
}

type ReadWrite16 interface {
	Read16
	Write16
}

type ReadWrite interface {
	ReadWrite8
	ReadWrite16
}

type ReadSafeWrite interface {
	ReadWrite
	ReadSafe
}

type HasCycles interface {
	GetCycles() uint
}
