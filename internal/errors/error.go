package errors

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/token"
)

type Error struct {
	Line  int
	Where string
	Err   error
}

func New(tok *token.Token, err error) *Error {
	e := &Error{
		Line:  tok.Line,
		Where: " at '" + tok.Lexeme + "'",
		Err:   err,
	}
	if tok.Type == token.TypeEOF {
		e.Where = " at end"
	}
	return e
}

func (err *Error) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", err.Line, err.Where, err.Err.Error())
}

func (err *Error) Unwrap() error {
	return err.Err
}
