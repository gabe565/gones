package gones

import (
	"bytes"
	"syscall/js"

	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/console"
	log "github.com/sirupsen/logrus"
)

func newConsole(_ string) (*console.Console, error) {
	jsData := js.Global().Get("cartridge")
	goData := make([]byte, jsData.Get("length").Int())
	js.CopyBytesToGo(goData, jsData)

	r := bytes.NewReader(goData)

	cart, err := cartridge.FromiNes(r)
	if err != nil {
		return nil, err
	}
	log.WithField("title", cart.Name()).Info("Loaded cartridge")

	js.Global().Call("SetRomName", cart.Name())

	c, err := console.New(cart)
	if err != nil {
		return c, err
	}

	js.Global().Set("Gones", js.ValueOf(map[string]any{
		"exit": js.FuncOf(func(this js.Value, args []js.Value) any {
			c.SetUpdateAction(console.ActionExit)
			return nil
		}),
		"saveState": js.FuncOf(func(this js.Value, args []js.Value) any {
			c.SetUpdateAction(console.ActionSaveState)
			return nil
		}),
		"loadState": js.FuncOf(func(this js.Value, args []js.Value) any {
			c.SetUpdateAction(console.ActionLoadState)
			return nil
		}),
	}))

	return c, nil
}
