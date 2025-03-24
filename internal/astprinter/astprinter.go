package astprinter

import (
	"strings"

	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/loxtype"
)

func Print(e expr.Expr) string {
	return ASTPrinter{}.Print(e).String()
}

type ASTPrinter struct{}

var _ expr.Visitor = ASTPrinter{}

func (ap ASTPrinter) Print(e expr.Expr) loxtype.Type {
	return e.Accept(ap)
}

func (ap ASTPrinter) parenthesize(name string, expressions ...expr.Expr) loxtype.Type {
	var builder strings.Builder
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, e := range expressions {
		builder.WriteRune(' ')
		builder.WriteString(e.Accept(ap).String())
	}
	builder.WriteRune(')')
	return loxtype.String(builder.String())
}

func (ap ASTPrinter) VisitBinary(e *expr.Binary) loxtype.Type {
	return ap.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (ap ASTPrinter) VisitGrouping(e *expr.Grouping) loxtype.Type {
	return ap.parenthesize("group", e.Expression)
}

func (ap ASTPrinter) VisitLiteral(e *expr.Literal) loxtype.Type {
	if e.Value == nil {
		return loxtype.Nil{}
	}
	return e.Value
}

func (ap ASTPrinter) VisitUnary(e *expr.Unary) loxtype.Type {
	return ap.parenthesize(e.Operator.Lexeme, e.Right)
}
