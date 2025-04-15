package interpreter

import (
	"errors"
	"fmt"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
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
	Call(*environment.Environment, *Interpreter, []loxtype.Type) (loxtype.Type, error)
}

type function ast.FunctionStmt

var _ callable = (*function)(nil)

func (f *function) Arity() int                        { return len(f.Params) }
func (*function) Equals(loxtype.Type) loxtype.Boolean { return false }
func (*function) IsTruthy() loxtype.Boolean           { return true }
func (f *function) String() string                    { return fmt.Sprintf("<fn: %s>", f.Name.Lexeme) }

func (f *function) Call(parent *environment.Environment, i *Interpreter, args []loxtype.Type) (loxtype.Type, error) {
	env := parent.MakeChild()
	for i, param := range f.Params {
		env.Define(param, args[i])
	}
	var val *returnValue
	if err := i.executeBlock(env, f.Body); err != nil && !errors.As(err, &val) {
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

func (nf *nativeFunction) Call(_ *environment.Environment, i *Interpreter, args []loxtype.Type) (loxtype.Type, error) {
	return nf.impl(i, args)
}
