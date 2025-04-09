package parser_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
)

func TestParser_Parse(t *testing.T) {
	code := `1 + "hello" * (3 - 4) > 14 == true;`
	tokens, scanErr := scanner.New(code).ScanTokens()
	require.NoError(t, scanErr)
	stmts, parseErr := parser.New(tokens).Parse()
	require.NoError(t, parseErr)
	var p ast.Printer

	for _, s := range stmts {
		t.Log(p.Print(s))
	}
}
