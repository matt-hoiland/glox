package expr

import "strings"

type ASTPrinter struct{}

var _ Visitor[string] = ASTPrinter{}

func (ap ASTPrinter) Print(expr Expr[string]) string {
	return expr.Accept(ap)
}

func (ap ASTPrinter) parenthesize(name string, exprs ...Expr[string]) string {
	var builder strings.Builder
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteRune(' ')
		builder.WriteString(expr.Accept(ap))
	}
	builder.WriteRune(')')
	return builder.String()
}

func (ap ASTPrinter) VisitBinary(expr *Binary[string]) string {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (ap ASTPrinter) VisitGrouping(expr *Grouping[string]) string {
	return ap.parenthesize("group", expr.Expression)
}

func (ap ASTPrinter) VisitLiteral(expr *Literal[string]) string {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value.String()
}

func (ap ASTPrinter) VisitUnary(expr *Unary[string]) string {
	return ap.parenthesize(expr.Operator.Lexeme, expr.Right)
}
