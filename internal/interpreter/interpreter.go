package interpreter

import (
	"errors"
	"fmt"
	"io"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/matt-hoiland/glox/internal/token"
)

var (
	ErrUnimplemented = errors.New("unimplemented")

	ErrType           = errors.New("type-error")
	ErrNonBooleanType = fmt.Errorf("non-boolean %w", ErrType)
	ErrNonNumericType = fmt.Errorf("non-numeric %w", ErrType)
	ErrNonStringType  = fmt.Errorf("non-string %w", ErrType)
)

type Interpreter struct {
	w io.Writer
}

func New(w io.Writer) *Interpreter {
	return &Interpreter{w}
}

func (i *Interpreter) Run(env *environment.Environment, code string) error {
	var (
		tokens []*token.Token
		stmts  []ast.Stmt
		err    error
	)

	if tokens, err = scanner.New(code).ScanTokens(); err != nil {
		return err
	}

	if stmts, err = parser.New(tokens).Parse(); err != nil {
		return err
	}

	if err = i.Interpret(env, stmts); err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) Interpret(env *environment.Environment, stmts []ast.Stmt) error {
	for _, s := range stmts {
		if err := i.execute(env, s); err != nil {
			return err
		}
	}
	return nil
}
