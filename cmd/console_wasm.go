package cmd

import (
	"bytes"
	"syscall/js"

	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/console"
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

	return console.New(cart)
}
