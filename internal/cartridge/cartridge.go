package cartridge

import (
	"crypto/md5"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/database"
	"gabe565.com/gones/internal/interrupt"
)

type Cartridge struct {
	hash   string
	name   string
	Header INESFileHeader `msgpack:"-"`

	PRG     []byte `msgpack:"-"`
	CHR     []byte `msgpack:"alias:Chr"`
	SRAM    []byte `msgpack:"alias:Sram"`
	Mirror  Mirror
	Battery bool `msgpack:"-"`
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

	cart.PRG = make([]byte, consts.PRGROMAddr, consts.PRGChunkSize*2)
	cart.PRG = append(cart.PRG, b...)
	cart.PRG = cart.PRG[:cap(cart.PRG)]
	cart.PRG[interrupt.ResetVector+1-consts.PRGChunkSize*2] = 0x86

	cart.CHR = make([]byte, consts.CHRChunkSize)

	return cart
}

func (c *Cartridge) Name() string {
	return c.name
}

func (c *Cartridge) SetName(path string) {
	c.name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func (c *Cartridge) Hash() string {
	return c.hash
}

func (c *Cartridge) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("title", c.name),
	)
}
