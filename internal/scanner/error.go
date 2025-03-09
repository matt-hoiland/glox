package scanner

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpectedRune     = errors.New("unexpected rune")
	ErrUnterminatedString = errors.New("unterminated string")
)

type Error struct {
	Line  int
	Where string
	Err   error
}

func (err *Error) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s\n", err.Line, err.Where, err.Err.Error())
}

func (err *Error) Unwrap() error {
	return err.Err
}
