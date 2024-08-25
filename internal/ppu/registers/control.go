package registers

// PPU Control bits
//
// 7 6 5 4 3 2 1 0
// V P H B S I N N
// ╷ ╷ ╷ ╷ ╷ ╷ ╷ ╷
// │ │ │ │ │ │ └─┴╴Base nametable address
// │ │ │ │ │ │       (0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00)
// │ │ │ │ │ └────╴VRAM address increment per CPU read/write of PPUDATA
// │ │ │ │ │          (0: add 1, going across; 1: add 32, going down)
// │ │ │ │ └──────╴Sprite pattern table address for 8x8 sprites
// │ │ │ │            (0: $0000; 1: $1000; ignored in 8x16 mode)
// │ │ │ └────────╴Background pattern table address (0: $0000; 1: $1000)
// │ │ └──────────╴Sprite size (0: 8x8 pixels; 1: 8x16 pixels – see PPU OAM#Byte 1)
// │ └────────────╴PPU master/slave select
// │                  (0: read backdrop from EXT pins; 1: output color on EXT pins)
// └──────────────╴Generate an NMI at the start of the vertical blanking interval
//                    (0: off; 1: on)

type Control struct {
	NametableX        bool
	NametableY        bool
	IncrementMode     bool
	SpriteTileSelect  bool
	BgTileSelect      bool
	SpriteHeight      bool
	MasterSlaveSelect bool
	EnableNMI         bool
}

const (
	CtrlNametableX = 1 << iota
	CtrlNametableY
	CtrlIncrementMode
	CtrlSpriteTileSelect
	CtrlBgTileSelect
	CtrlSpriteHeight
	CtrlMasterSlaveSelect
	CtrlEnableNMI
)

func (c *Control) Set(data byte) {
	c.NametableX = data&CtrlNametableX != 0
	c.NametableY = data&CtrlNametableY != 0
	c.IncrementMode = data&CtrlIncrementMode != 0
	c.SpriteTileSelect = data&CtrlSpriteTileSelect != 0
	c.BgTileSelect = data&CtrlBgTileSelect != 0
	c.SpriteHeight = data&CtrlSpriteHeight != 0
	c.MasterSlaveSelect = data&CtrlMasterSlaveSelect != 0
	c.EnableNMI = data&CtrlEnableNMI != 0
}

func (c Control) VRAMAddr() byte {
	if c.IncrementMode {
		return 32
	}
	return 1
}

func (c Control) SpriteTileAddr() uint16 {
	if c.SpriteTileSelect {
		return 0x1000
	}
	return 0
}

func (c Control) BgTileAddr() uint16 {
	if c.BgTileSelect {
		return 0x1000
	}
	return 0
}

func (c Control) SpriteSize() byte {
	if c.SpriteHeight {
		return 16
	}
	return 8
}
