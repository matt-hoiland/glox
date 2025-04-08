package interpreter

import (
	"errors"
	"fmt"

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
	environment *environment.Environment
}

var (
	_ ast.ExprVisitor = (*Interpreter)(nil)
	_ ast.StmtVisitor = (*Interpreter)(nil)
)

func New() *Interpreter {
	return &Interpreter{
		environment: environment.New(),
	}
}

func (i *Interpreter) Run(code string) error {
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

	if err = i.Interpret(stmts); err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) error {
	for _, s := range stmts {
		if err := i.execute(s); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) execute(s ast.Stmt) error {
	if _, err := s.Accept(i); err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) evaluate(e ast.Expr) (loxtype.Type, error) {
	return e.Accept(i)
}

func (i *Interpreter) isEqual(a, b loxtype.Type) loxtype.Boolean {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	ae, ok := a.(loxtype.Equalser)
	if !ok {
		return false
	}
	return ae.Equals(b)
}

func (i *Interpreter) isTruthy(value loxtype.Type) loxtype.Boolean {
	if value == nil {
		return false
	}
	if _, ok := value.(loxtype.Nil); ok {
		return false
	}
	if b, ok := value.(loxtype.Boolean); ok {
		return b
	}
	return true
}

func (i *Interpreter) VisitExpressionStmt(s *ast.ExpressionStmt) (loxtype.Type, error) {
	_, err := i.evaluate(s.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(s *ast.PrintStmt) (loxtype.Type, error) {
	value, err := i.evaluate(s.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(value)
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(s *ast.VarStmt) (loxtype.Type, error) {
	var (
		value loxtype.Type = loxtype.Nil{}
		err   error
	)
	if s.Initializer != nil {
		if value, err = i.evaluate(s.Initializer); err != nil {
			return nil, err
		}
	}

	i.environment.Define(s.Name, value)
	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(e *ast.AssignExpr) (loxtype.Type, error) {
	value, err := i.evaluate(e.Value)
	if err != nil {
		return nil, err
	}
	if err = i.environment.Assign(e.Name, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitBinaryExpr(e *ast.BinaryExpr) (loxtype.Type, error) {
	left, err := i.evaluate(e.Left)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate left operand of binary expression: %w", err)
	}
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate right operand of binary expression: %w", err)
	}

	switch e.Operator.Type {
	case token.TypeBangEqual:
		return !i.isEqual(left, right), nil
	case token.TypeEqualEqual:
		return i.isEqual(left, right), nil
	case token.TypeGreater:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("greater expression: %w", err)
		}
		return a.Greater(b), nil
	case token.TypeGreaterEqual:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("greater-equal expression: %w", err)
		}
		return a.GreaterEqual(b), nil
	case token.TypeLess:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("less expression: %w", err)
		}
		return a.Less(b), nil
	case token.TypeLessEqual:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("less-equal expression: %w", err)
		}
		return a.LessEqual(b), nil
	case token.TypeMinus:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("minus expression: %w", err)
		}
		return a.Subtract(b), nil
	case token.TypeSlash:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("slash expression: %w", err)
		}
		return a.Divide(b), nil
	case token.TypeStar:
		a, b, err := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		if err != nil {
			return nil, fmt.Errorf("star expression: %w", err)
		}
		return a.Multiply(b), nil
	case token.TypePlus:
		na, nb, nErr := convertBoth[loxtype.Number](left, right, ErrNonNumericType)
		sa, sb, sErr := convertBoth[loxtype.String](left, right, ErrNonStringType)
		if nErr != nil && sErr != nil {
			return nil, fmt.Errorf("operands to plus expression must be either string or numeric: %w", ErrType)
		}
		if nErr != nil {
			return sa.Add(sb), nil
		}
		return na.Add(nb), nil
	}

	return nil, ErrUnimplemented
}

func (i *Interpreter) VisitGroupingExpr(e *ast.GroupingExpr) (loxtype.Type, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitLiteralExpr(e *ast.LiteralExpr) (loxtype.Type, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(e *ast.UnaryExpr) (loxtype.Type, error) {
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate operand of unary expression: %w", err)
	}

	switch e.Operator.Type {
	case token.TypeBang:
		return i.isTruthy(right).Negate(), nil
	case token.TypeMinus:
		n, ok := right.(loxtype.Number)
		if !ok {
			return nil, fmt.Errorf("cannot apply minus operator: %w", ErrNonNumericType)
		}
		return n.Negate(), nil
	}

	return nil, ErrUnimplemented
}

func (i *Interpreter) VisitVariableExpr(e *ast.VariableExpr) (loxtype.Type, error) {
	return i.environment.Get(e.Name)
}

func convertBoth[T any](a, b loxtype.Type, err error) (T, T, error) {
	at, aok := a.(T)
	bt, bok := b.(T)
	if !aok {
		return at, bt, fmt.Errorf("type-error: left-hand operand: %w", err)
	}
	if !bok {
		return at, bt, fmt.Errorf("type-error: right-hand operand: %w", err)
	}
	return at, bt, nil
}
