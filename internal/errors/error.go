package errors

import (
	"fmt"
)

type Error struct {
	Line  int
	Where string
	Err   error
}

func (err *Error) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", err.Line, err.Where, err.Err.Error())
}

func (err *Error) Unwrap() error {
	return err.Err
}
