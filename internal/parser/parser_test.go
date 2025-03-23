package parser_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
)

func TestParser_Parse(t *testing.T) {
	code := `1 + "hello" * (3 - 4) > 14 == true;`
	tokens, scanErr := scanner.New(code).ScanTokens()
	require.NoError(t, scanErr)
	ast, parseErr := parser.New(tokens).Parse()
	require.NoError(t, parseErr)
	var printer expr.ASTPrinter
	fmt.Println(printer.Print(ast))
}
