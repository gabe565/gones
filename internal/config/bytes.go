package config

import (
	"strings"

	"gabe565.com/utils/bytefmt"
)

type Bytes int64

func (b Bytes) MarshalText() ([]byte, error) {
	formatted := bytefmt.Encode(int64(b))
	formatted = strings.Replace(formatted, ".00", "", 1)
	return []byte(formatted), nil
}

func (b *Bytes) UnmarshalText(text []byte) error {
	v, err := bytefmt.Decode(string(text))
	if err != nil {
		return err
	}
	*b = Bytes(v)
	return nil
}
