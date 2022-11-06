package callbacks

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/cpu"
)

type Callback func(*cpu.CPU) error
type CallbackHandler func(*pixelgl.Window) Callback

var Callbacks = map[string]CallbackHandler{
	"snake": Snake,
}
