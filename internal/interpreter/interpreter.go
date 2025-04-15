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

type Interpreter struct {
	w       io.Writer
	globals *environment.Environment
	locals  map[ast.Expr]int
}

func New(w io.Writer) *Interpreter {
	env := &Interpreter{
		w:       w,
		globals: environment.New(),
		locals:  map[ast.Expr]int{},
	}

	env.globals.Define(
		&token.Token{
			Type:   token.TypeIdentifier,
			Lexeme: "clock",
		},
		&nativeFunction{
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
		&nativeFunction{
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

	if err = newResolver(i).resolveStmts(stmts); err != nil {
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

func (i *Interpreter) lookUpVariable(env *environment.Environment, name *token.Token,
	expr *ast.VariableExpr) (loxtype.Type, error) {
	distance := i.locals[expr]
	return env.GetAt(name, distance)
}

func (i *Interpreter) resolve(expr ast.Expr, distance int) {
	i.locals[expr] = distance
}
