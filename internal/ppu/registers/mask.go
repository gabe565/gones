package registers

// PPU Mask bits
//
// 7 6 5 4 3 2 1 0
// B G R s b M m G
// ╷ ╷ ╷ ╷ ╷ ╷ ╷ ╷
// │ │ │ │ │ │ │ └╴Greyscale (0: normal color, 1: produce a greyscale display)
// │ │ │ │ │ │ └──╴1: Show background in leftmost 8 pixels of screen, 0: Hide
// │ │ │ │ │ └────╴1: Show sprites in leftmost 8 pixels of screen, 0: Hide
// │ │ │ │ └──────╴1: Show background
// │ │ │ └────────╴1: Show sprites
// │ │ └──────────╴Emphasize red (green on PAL/Dendy)
// │ └────────────╴Emphasize green (red on PAL/Dendy)
// └──────────────╴Emphasize blue

type Mask struct {
	Grayscale           bool
	BgLeftColEnable     bool
	SpriteLeftColEnable bool
	BackgroundEnable    bool
	SpriteEnable        bool
	EmphasizeRed        bool
	EmphasizeGreen      bool
	EmphasizeBlue       bool
}

const (
	MaskGrayscale = 1 << iota
	MaskBgLeftColEnable
	MaskSpriteLeftColEnable
	MaskBackgroundEnable
	MaskSpriteEnable
	MaskEmphasizeRed
	MaskEmphasizeGreen
	MaskEmphasizeBlue
)

func (m *Mask) Set(data byte) {
	m.Grayscale = data&MaskGrayscale != 0
	m.BgLeftColEnable = data&MaskBgLeftColEnable != 0
	m.SpriteLeftColEnable = data&MaskSpriteLeftColEnable != 0
	m.BackgroundEnable = data&MaskBackgroundEnable != 0
	m.SpriteEnable = data&MaskSpriteEnable != 0
	m.EmphasizeRed = data&MaskEmphasizeRed != 0
	m.EmphasizeGreen = data&MaskEmphasizeGreen != 0
	m.EmphasizeBlue = data&MaskEmphasizeBlue != 0
}

func (m *Mask) Get() byte {
	var data byte
	if m.Grayscale {
		data |= MaskGrayscale
	}
	if m.BgLeftColEnable {
		data |= MaskBgLeftColEnable
	}
	if m.SpriteLeftColEnable {
		data |= MaskSpriteLeftColEnable
	}
	if m.BackgroundEnable {
		data |= MaskBackgroundEnable
	}
	if m.SpriteEnable {
		data |= MaskSpriteEnable
	}
	if m.EmphasizeRed {
		data |= MaskEmphasizeRed
	}
	if m.EmphasizeGreen {
		data |= MaskEmphasizeGreen
	}
	if m.EmphasizeBlue {
		data |= MaskEmphasizeBlue
	}
	return data
}

func (m *Mask) RenderingEnabled() bool {
	return m.BackgroundEnable || m.SpriteEnable
}
