package ast_test

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

func ExamplePrinter() {
	expression := ast.NewBinaryExpr(
		ast.NewUnaryExpr(
			token.NewToken(token.TypeMinus, "-", nil, 1),
			ast.NewLiteralExpr(loxtype.Number(123)),
		),
		token.NewToken(token.TypeStar, "*", nil, 1),
		ast.NewGroupingExpr(
			ast.NewLiteralExpr(nil),
		),
	)

	var p ast.Printer
	s, _ := p.Print(expression)
	fmt.Println(s)
	// Output: (* (- 123) (group nil))
}
