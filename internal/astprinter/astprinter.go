package astprinter

import (
	"strings"

	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/loxtype"
)

func Print(e expr.Expr) string {
	s, _ := ASTPrinter{}.Print(e)
	return s.String()
}

type ASTPrinter struct{}

var _ expr.Visitor = ASTPrinter{}

func (ap ASTPrinter) Print(e expr.Expr) (loxtype.Type, error) {
	return e.Accept(ap)
}

func (ap ASTPrinter) parenthesize(name string, expressions ...expr.Expr) (loxtype.Type, error) {
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

func (ap ASTPrinter) VisitBinary(e *expr.Binary) (loxtype.Type, error) {
	return ap.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (ap ASTPrinter) VisitGrouping(e *expr.Grouping) (loxtype.Type, error) {
	return ap.parenthesize("group", e.Expression)
}

func (ap ASTPrinter) VisitLiteral(e *expr.Literal) (loxtype.Type, error) {
	if e.Value == nil {
		return loxtype.Nil{}, nil
	}
	return e.Value, nil
}

func (ap ASTPrinter) VisitUnary(e *expr.Unary) (loxtype.Type, error) {
	return ap.parenthesize(e.Operator.Lexeme, e.Right)
}
