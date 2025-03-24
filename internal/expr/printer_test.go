package expr_test

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

func ExampleASTPrinter() {
	expression := expr.NewBinary(
		expr.NewUnary(
			token.NewToken(token.TypeMinus, "-", nil, 1),
			expr.NewLiteral[string](loxtype.Number(123)),
		),
		token.NewToken(token.TypeStar, "*", nil, 1),
		expr.NewGrouping(
			expr.NewLiteral[string](nil),
		),
	)

	var printer expr.ASTPrinter
	fmt.Println(printer.Print(expression))
	// Output: (* (- 123) (group nil))
}
