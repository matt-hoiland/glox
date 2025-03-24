package astprinter_test

import (
	"fmt"

	"github.com/matt-hoiland/glox/internal/astprinter"
	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/token"
)

func ExampleASTPrinter() {
	expression := expr.NewBinary(
		expr.NewUnary(
			token.NewToken(token.TypeMinus, "-", nil, 1),
			expr.NewLiteral(loxtype.Number(123)),
		),
		token.NewToken(token.TypeStar, "*", nil, 1),
		expr.NewGrouping(
			expr.NewLiteral(nil),
		),
	)

	var p astprinter.ASTPrinter
	s, _ := p.Print(expression)
	fmt.Println(s)
	// Output: (* (- 123) (group nil))
}
