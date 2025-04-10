package interpreter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
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

type Callable interface {
	loxtype.Type
	Arity() int
	Call(*Interpreter, []loxtype.Type) (loxtype.Type, error)
}

type NativeFunction struct {
	name  string
	arity int
	impl  func(*Interpreter, []loxtype.Type) (loxtype.Type, error)
}

var _ Callable = (*NativeFunction)(nil)

func (nf *NativeFunction) Arity() int { return nf.arity }

func (nf *NativeFunction) Call(i *Interpreter, args []loxtype.Type) (loxtype.Type, error) {
	return nf.impl(i, args)
}

func (*NativeFunction) Equals(loxtype.Type) loxtype.Boolean { return false }
func (*NativeFunction) IsTruthy() loxtype.Boolean           { return true }
func (nf *NativeFunction) String() string                   { return fmt.Sprintf("<native fn: %s>", nf.name) }

type Interpreter struct {
	w       io.Writer
	globals *environment.Environment
}

func New(w io.Writer) *Interpreter {
	env := &Interpreter{
		w:       w,
		globals: environment.New(),
	}

	env.globals.Define(
		&token.Token{
			Type:   token.TypeIdentifier,
			Lexeme: "clock",
		},
		&NativeFunction{
			name:  "clock",
			arity: 0,
			impl: func(*Interpreter, []loxtype.Type) (loxtype.Type, error) {
				return loxtype.Number(time.Now().UnixMilli()), nil
			},
		},
	)

	env.globals.Define(
		&token.Token{
			Type:   token.TypeIdentifier,
			Lexeme: "exit",
		},
		&NativeFunction{
			name:  "exit",
			arity: 0,
			impl: func(*Interpreter, []loxtype.Type) (loxtype.Type, error) {
				os.Exit(0)
				return nil, nil //nolint: nilnil // native function
			},
		},
	)

	return env
}

func (i *Interpreter) Run(code string) error {
	var (
		env    = i.globals.MakeChild()
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

	if err = i.interpret(env, stmts); err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) Evaluate(expr ast.Expr) (loxtype.Type, error) {
	return i.evaluate(i.globals, expr)
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) error {
	return i.interpret(i.globals, stmts)
}

func (i *Interpreter) interpret(env *environment.Environment, stmts []ast.Stmt) error {
	for _, s := range stmts {
		if err := i.execute(env, s); err != nil {
			return err
		}
	}
	return nil
}
