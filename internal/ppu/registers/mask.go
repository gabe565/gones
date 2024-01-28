package registers

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
