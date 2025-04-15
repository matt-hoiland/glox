package interpreter

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

var _ ast.ExprVisitor = (*Interpreter)(nil)

func (i *Interpreter) evaluate(env *environment.Environment, e ast.Expr) (loxtype.Type, error) {
	return e.Accept(env, i)
}

func (i *Interpreter) VisitAssignExpr(env *environment.Environment, e *ast.AssignExpr) (loxtype.Type, error) {
	value, err := i.evaluate(env, e.Value)
	if err != nil {
		return nil, err
	}
	if err = env.Assign(e.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitBinaryExpr(env *environment.Environment, e *ast.BinaryExpr) (loxtype.Type, error) {
	left, err := i.evaluate(env, e.Left)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate left operand of binary expression: %w", err)
	}
	right, err := i.evaluate(env, e.Right)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate right operand of binary expression: %w", err)
	}

	switch e.Operator.Type {
	case token.TypeBangEqual:
		return !left.Equals(right), nil
	case token.TypeEqualEqual:
		return left.Equals(right), nil
	case token.TypeGreater:
		return i.biopGreater(left, right)
	case token.TypeGreaterEqual:
		return i.biopGreaterEqual(left, right)
	case token.TypeLess:
		return i.biopLess(left, right)
	case token.TypeLessEqual:
		return i.biopLessEqual(left, right)
	case token.TypeMinus:
		return i.biopMinus(left, right)
	case token.TypeSlash:
		return i.biopSlash(left, right)
	case token.TypeStar:
		return i.biopStar(left, right)
	case token.TypePlus:
		return i.biopPlus(left, right)
	default:
		return nil, ErrUnimplemented
	}
}

func (*Interpreter) biopGreater(left, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("greater expression: %w", ErrNonNumericType)
	}
	return a.Greater(b), nil
}

func (*Interpreter) biopGreaterEqual(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("greater-equal expression: %w", ErrNonNumericType)
	}
	return a.GreaterEqual(b), nil
}

func (*Interpreter) biopLess(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("less expression: %w", ErrNonNumericType)
	}
	return a.Less(b), nil
}

func (*Interpreter) biopLessEqual(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("less-equal expression: %w", ErrNonNumericType)
	}
	return a.LessEqual(b), nil
}

func (*Interpreter) biopMinus(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("minus expression: %w", ErrNonNumericType)
	}
	return a.Subtract(b), nil
}

func (*Interpreter) biopSlash(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("slash expression: %w", ErrNonNumericType)
	}
	return a.Divide(b), nil
}

func (*Interpreter) biopStar(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	a, b, ok := convertBoth[loxtype.Number](left, right)
	if !ok {
		return nil, fmt.Errorf("star expression: %w", ErrNonNumericType)
	}
	return a.Multiply(b), nil
}

func (*Interpreter) biopPlus(left loxtype.Type, right loxtype.Type) (loxtype.Type, error) {
	na, nb, nok := convertBoth[loxtype.Number](left, right)
	sa, sb, sok := convertBoth[loxtype.String](left, right)
	if !nok && !sok {
		return nil, fmt.Errorf("operands to plus expression must be either string or numeric: %w", ErrType)
	}
	if !nok {
		return sa.Add(sb), nil
	}
	return na.Add(nb), nil
}

func convertBoth[T any](a, b loxtype.Type) (T, T, bool) {
	at, aok := a.(T)
	bt, bok := b.(T)
	if !aok {
		return at, bt, false
	}
	if !bok {
		return at, bt, false
	}
	return at, bt, true
}

func (i *Interpreter) VisitCallExpr(env *environment.Environment, e *ast.CallExpr) (loxtype.Type, error) {
	var (
		callee    loxtype.Type
		arguments []loxtype.Type
		function  callable
		ok        bool
		err       error
	)

	if callee, err = i.evaluate(env, e.Callee); err != nil {
		return nil, err
	}

	for _, argExpr := range e.Arguments {
		var arg loxtype.Type
		if arg, err = i.evaluate(env, argExpr); err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	if function, ok = callee.(callable); !ok {
		return nil, fmt.Errorf("%s not a callable: can only call functions and classes", function)
	}

	if len(arguments) != function.Arity() {
		return nil, fmt.Errorf("expected %d arguments but got %d", function.Arity(), len(arguments))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGroupingExpr(env *environment.Environment, e *ast.GroupingExpr) (loxtype.Type, error) {
	return i.evaluate(env, e.Expression)
}

func (i *Interpreter) VisitLiteralExpr(_ *environment.Environment, e *ast.LiteralExpr) (loxtype.Type, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitLogicalExpr(env *environment.Environment, e *ast.LogicalExpr) (loxtype.Type, error) {
	left, err := i.evaluate(env, e.Left)
	if err != nil {
		return nil, err
	}

	if e.Operator.Type == token.TypeOr {
		if left.IsTruthy() {
			return left, nil
		}
	} else {
		if !left.IsTruthy() {
			return left, nil
		}
	}

	return i.evaluate(env, e.Right)
}

func (i *Interpreter) VisitUnaryExpr(env *environment.Environment, e *ast.UnaryExpr) (loxtype.Type, error) {
	right, err := i.evaluate(env, e.Right)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate operand of unary expression: %w", err)
	}

	switch e.Operator.Type {
	case token.TypeBang:
		return right.IsTruthy().Negate(), nil
	case token.TypeMinus:
		n, ok := right.(loxtype.Number)
		if !ok {
			return nil, fmt.Errorf("cannot apply minus operator: %w", ErrNonNumericType)
		}
		return n.Negate(), nil
	default:
		return nil, ErrUnimplemented
	}
}

func (i *Interpreter) VisitVariableExpr(env *environment.Environment, e *ast.VariableExpr) (loxtype.Type, error) {
	return env.Get(e.Name)
}
