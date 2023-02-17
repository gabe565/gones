package console

import (
	"errors"
	"io"
)

type ConsoleTest struct {
	Console *Console
	ResetIn uint16

	Callback func(b *ConsoleTest) error
}

func NewConsoleTest(r io.ReadSeeker, callback func(console *ConsoleTest) error) (*ConsoleTest, error) {
	console, err := stubConsole(r)
	if err != nil {
		return nil, err
	}

	return &ConsoleTest{
		Console:  console,
		Callback: callback,
	}, nil
}

func (b *ConsoleTest) Run() error {
	for {
		if b.ResetIn != 0 {
			b.ResetIn -= 1
			if b.ResetIn == 0 {
				b.Console.Reset()
			}
		}

		if b.Callback != nil {
			if err := b.Callback(b); err != nil {
				if errors.Is(err, ErrExit) {
					return nil
				} else {
					return err
				}
			}
		}

		if err := b.Console.Step(); err != nil {
			return err
		}
	}
}
