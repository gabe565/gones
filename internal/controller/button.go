package controller

import (
	"errors"
	"fmt"
)

//go:generate stringer -type Button -linecomment

type Button uint8

const (
	ButtonA Button = iota // a
	ButtonB               // b
	Select                // select
	Start                 // start
	Up                    // up
	Down                  // down
	Left                  // left
	Right                 // right
)

var ErrInvalidButton = errors.New("invalid button")

func (i *Button) UnmarshalText(b []byte) error {
	s := string(b)
	for j := 0; j < len(_Button_index)-1; j++ {
		if s == _Button_name[_Button_index[j]:_Button_index[j+1]] {
			*i = Button(j)
			return nil
		}
	}
	return fmt.Errorf("%v: %s", ErrInvalidButton, s)
}
