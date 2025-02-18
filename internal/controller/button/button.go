package button

import (
	"errors"
	"fmt"
)

//go:generate go tool stringer -type Button -linecomment

type Button uint8

const (
	A      Button = iota // a
	B                    // b
	Select               // select
	Start                // start
	Up                   // up
	Down                 // down
	Left                 // left
	Right                // right
)

var ErrInvalidButton = errors.New("invalid button")

func (i *Button) UnmarshalText(b []byte) error {
	s := string(b)
	for j := range len(_Button_index) - 1 {
		if s == _Button_name[_Button_index[j]:_Button_index[j+1]] {
			*i = Button(j)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrInvalidButton, s)
}
