package registers

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
	} else {
		return 1
	}
}

func (c Control) SpriteTileAddr() uint16 {
	if c.SpriteTileSelect {
		return 0x1000
	} else {
		return 0
	}
}

func (c Control) BgTileAddr() uint16 {
	if c.BgTileSelect {
		return 0x1000
	} else {
		return 0
	}
}

func (c Control) SpriteSize() byte {
	if c.SpriteHeight {
		return 16
	} else {
		return 8
	}
}
