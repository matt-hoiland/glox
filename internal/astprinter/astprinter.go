package astprinter

import (
	"strings"

	"github.com/matt-hoiland/glox/internal/expr"
)

func Print(e expr.Expr[string]) string {
	return ASTPrinter{}.Print(e)
}

type ASTPrinter struct{}

var _ expr.Visitor[string] = ASTPrinter{}

func (ap ASTPrinter) Print(e expr.Expr[string]) string {
	return e.Accept(ap)
}

func (ap ASTPrinter) parenthesize(name string, expressions ...expr.Expr[string]) string {
	var builder strings.Builder
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, e := range expressions {
		builder.WriteRune(' ')
		builder.WriteString(e.Accept(ap))
	}
	builder.WriteRune(')')
	return builder.String()
}

func (ap ASTPrinter) VisitBinary(e *expr.Binary[string]) string {
	return ap.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (ap ASTPrinter) VisitGrouping(e *expr.Grouping[string]) string {
	return ap.parenthesize("group", e.Expression)
}

func (ap ASTPrinter) VisitLiteral(e *expr.Literal[string]) string {
	if e.Value == nil {
		return "nil"
	}
	return e.Value.String()
}

func (ap ASTPrinter) VisitUnary(e *expr.Unary[string]) string {
	return ap.parenthesize(e.Operator.Lexeme, e.Right)
}
