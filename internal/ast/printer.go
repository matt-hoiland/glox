package ast

import (
	"strings"

	"github.com/matt-hoiland/glox/internal/loxtype"
)

func Print(e Expr) string {
	s, _ := Printer{}.Print(e)
	return s.String()
}

type Printer struct{}

var _ ExprVisitor = Printer{}

func (ap Printer) Print(e Expr) (loxtype.Type, error) {
	return e.Accept(ap)
}

func (ap Printer) parenthesize(name string, expressions ...Expr) (loxtype.Type, error) {
	var builder strings.Builder
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, e := range expressions {
		builder.WriteRune(' ')
		s, _ := e.Accept(ap)
		builder.WriteString(s.String())
	}
	builder.WriteRune(')')
	return loxtype.String(builder.String()), nil
}

func (ap Printer) VisitBinaryExpr(e *BinaryExpr) (loxtype.Type, error) {
	return ap.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (ap Printer) VisitGroupingExpr(e *GroupingExpr) (loxtype.Type, error) {
	return ap.parenthesize("group", e.Expression)
}

func (ap Printer) VisitLiteralExpr(e *LiteralExpr) (loxtype.Type, error) {
	if e.Value == nil {
		return loxtype.Nil{}, nil
	}
	return e.Value, nil
}

func (ap Printer) VisitUnaryExpr(e *UnaryExpr) (loxtype.Type, error) {
	return ap.parenthesize(e.Operator.Lexeme, e.Right)
}
