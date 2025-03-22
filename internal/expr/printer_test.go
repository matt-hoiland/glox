package expr_test

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/matt-hoiland/glox/internal/scanner/literal"
	"github.com/matt-hoiland/glox/internal/scanner/tokentype"
)

func ExampleASTPrinter() {
	expression := expr.NewBinary(
		expr.NewUnary(
			scanner.NewToken(tokentype.Minus, "-", nil, 1),
			expr.NewLiteral[string](literal.Number(123)),
		),
		scanner.NewToken(tokentype.Star, "*", nil, 1),
		expr.NewGrouping(
			expr.NewLiteral[string](nil),
		),
	)

	var printer expr.ASTPrinter
	fmt.Println(printer.Print(expression))
	// Output: (* (- 123) (group nil))
}
