package cartridge

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/database"
	"github.com/gabe565/gones/internal/interrupt"
)

type Cartridge struct {
	hash string
	name string

	prg       []byte
	CHR       []byte `msgpack:"alias:Chr"`
	SRAM      []byte `msgpack:"alias:Sram"`
	Mapper    uint8  `msgpack:"-"`
	Submapper uint8  `msgpack:"-"`
	Mirror    Mirror
	Battery   bool `msgpack:"-"`
}

func New() *Cartridge {
	return &Cartridge{
		SRAM: make([]byte, 0x2000),
	}
}

func FromBytes(b []byte) *Cartridge {
	cart := New()
	cart.hash = fmt.Sprintf("%x", md5.Sum(b))
	if cart.hash != "" {
		cart.name, _ = database.FindNameByHash(cart.hash)
	}

	cart.prg = make([]byte, consts.PRGROMAddr, consts.PRGChunkSize*2)
	cart.prg = append(cart.prg, b...)
	cart.prg = cart.prg[:cap(cart.prg)]
	cart.prg[interrupt.ResetVector+1-consts.PRGChunkSize*2] = 0x86

	cart.CHR = make([]byte, consts.CHRChunkSize)

	return cart
}

func (c *Cartridge) Name() string {
	return c.name
}

func (c *Cartridge) SetName(path string) {
	c.name = strings.TrimSuffix(filepath.Base(path), ".nes")
}

func (c *Cartridge) Hash() string {
	return c.hash
}

func (c *Cartridge) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", c.name),
	)
}
