package interpreter

import (
	"errors"
	"fmt"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
)

type functionType string

const (
	none     functionType = "NONE"
	function functionType = "FUNCTION"
)

//nolint:errname // Intentional abuse of go's error system.
type returnValue struct {
	value loxtype.Type
}

var _ error = (*returnValue)(nil)

func (r *returnValue) Error() string {
	return "return value"
}

func (i *Interpreter) VisitReturnStmt(env *environment.Environment, s *ast.ReturnStmt) (loxtype.Type, error) {
	var (
		value loxtype.Type
		err   error
	)

	if s.Value != nil {
		if value, err = i.evaluate(env, s.Value); err != nil {
			return nil, err
		}
	}

	return nil, &returnValue{value}
}

type callable interface {
	loxtype.Type
	Arity() int
	Call(*Interpreter, []loxtype.Type) (loxtype.Type, error)
}

type loxFunction struct {
	closure *environment.Environment
	stmt    *ast.FunctionStmt
}

var _ callable = (*loxFunction)(nil)

func newFunction(env *environment.Environment, stmt *ast.FunctionStmt) *loxFunction {
	return &loxFunction{
		closure: env,
		stmt:    stmt,
	}
}

func (f *loxFunction) Arity() int                        { return len(f.stmt.Params) }
func (*loxFunction) Equals(loxtype.Type) loxtype.Boolean { return false }
func (*loxFunction) IsTruthy() loxtype.Boolean           { return true }
func (f *loxFunction) String() string                    { return fmt.Sprintf("<fn: %s>", f.stmt.Name.Lexeme) }

func (f *loxFunction) Call(i *Interpreter, args []loxtype.Type) (loxtype.Type, error) {
	env := f.closure.MakeChild()
	for i, param := range f.stmt.Params {
		env.Define(param, args[i])
	}
	var val *returnValue
	if err := i.executeBlock(env, f.stmt.Body); err != nil && !errors.As(err, &val) {
		return nil, err
	}
	if val == nil || val.value == nil {
		return loxtype.Nil{}, nil
	}
	return val.value, nil
}

type nativeFunction struct {
	name  string
	arity int
	impl  func(*Interpreter, []loxtype.Type) (loxtype.Type, error)
}

var _ callable = (*nativeFunction)(nil)

func (nf *nativeFunction) Arity() int                       { return nf.arity }
func (*nativeFunction) Equals(loxtype.Type) loxtype.Boolean { return false }
func (*nativeFunction) IsTruthy() loxtype.Boolean           { return true }
func (nf *nativeFunction) String() string                   { return fmt.Sprintf("<native fn: %s>", nf.name) }

func (nf *nativeFunction) Call(i *Interpreter, args []loxtype.Type) (loxtype.Type, error) {
	return nf.impl(i, args)
}
