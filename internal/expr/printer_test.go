package expr_test

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/token"
	"github.com/matt-hoiland/glox/internal/token/tokentype"
)

func ExampleASTPrinter() {
	expression := expr.NewBinary(
		expr.NewUnary(
			token.NewToken(tokentype.Minus, "-", nil, 1),
			expr.NewLiteral[string](literal.Number(123)),
		),
		token.NewToken(tokentype.Star, "*", nil, 1),
		expr.NewGrouping(
			expr.NewLiteral[string](nil),
		),
	)

	var printer expr.ASTPrinter
	fmt.Println(printer.Print(expression))
	// Output: (* (- 123) (group nil))
}
