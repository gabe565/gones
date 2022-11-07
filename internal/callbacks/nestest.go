package callbacks

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gabe565/gones/internal/cpu"
	"sync"
)

func NesTest(win *pixelgl.Window) Callback {
	var once sync.Once

	return func(c *cpu.CPU) error {
		once.Do(func() {
			c.ProgramCounter = 0xC000
		})

		fmt.Println(c.Trace())
		return nil
	}
}
