package config

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Key ebiten.Key

func (k Key) MarshalText() ([]byte, error) {
	return ebiten.Key(k).MarshalText()
}

func (k *Key) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*k = Key(-1)
		return nil
	}

	var temp ebiten.Key
	if err := temp.UnmarshalText(text); err != nil {
		return err
	}
	*k = Key(temp)
	return nil
}
