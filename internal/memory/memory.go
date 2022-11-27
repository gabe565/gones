package memory

type Read interface {
	ReadMem(uint16) byte
}

type Write interface {
	WriteMem(uint16, byte)
}

type ReadWrite interface {
	Read
	Write
}
