package gones

import (
	"bytes"
	"log/slog"
	"os"
	"syscall/js"

	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/config"
	"gabe565.com/gones/internal/console"
	"github.com/hajimehoshi/ebiten/v2"
)

func newConsole(conf *config.Config, _ string) (*console.Console, error) {
	jsCartridge := js.Global().Get("GonesCartridge")
	jsData := jsCartridge.Get("data")
	goData := make([]byte, jsData.Get("length").Int())
	js.CopyBytesToGo(goData, jsData)

	r := bytes.NewReader(goData)

	cart, err := cartridge.FromiNes(r)
	if err != nil {
		return nil, err
	}
	if cart.Name() == "" {
		cart.SetName(jsCartridge.Get("name").String())
	}
	slog.Info("Loaded cartridge", "", cart)

	js.Global().Get("GonesClient").Call("setRomName", cart.Name())

	c, err := console.New(conf, cart)
	if err != nil {
		return c, err
	}

	js.Global().Set("Gones", js.Global().Get("Object").Call("freeze", js.ValueOf(map[string]any{
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
	})))

	_ = os.Setenv("EBITENGINE_SCREENSHOT_KEY", ebiten.Key(conf.Input.Screenshot).String())

	return c, nil
}
