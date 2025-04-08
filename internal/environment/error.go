package environment

import (
	"errors"
	"fmt"

	"github.com/matt-hoiland/glox/internal/token"
)

var ErrUndefinedVariable = errors.New("undefined variable")

type UndefinedVariableError struct {
	token *token.Token
}

var (
	_ error                       = (*UndefinedVariableError)(nil)
	_ interface{ Unwrap() error } = (*UndefinedVariableError)(nil)
)

func newUndefinedVariableError(token *token.Token) *UndefinedVariableError {
	return &UndefinedVariableError{token: token}
}

func (e *UndefinedVariableError) Error() string {
	return fmt.Sprintf("%s: %s", ErrUndefinedVariable.Error(), e.token.Lexeme)
}

func (*UndefinedVariableError) Unwrap() error {
	return ErrUndefinedVariable
}
