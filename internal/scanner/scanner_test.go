package scanner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/runes"
	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/matt-hoiland/glox/internal/scanner/tokentype"
)

func TestScanner_ScanTokens(t *testing.T) {
	t.Parallel()

	type Test struct {
		Name   string
		Source string
		Tokens []*scanner.Token
		Err    error
	}

	tests := []Test{
		{
			Name:   "success/empty source",
			Source: ``,
			Tokens: []*scanner.Token{
				{Type: tokentype.EOF},
			},
		},
		{
			Name:   "success/assign string to variable",
			Source: `var language = "lox";`,
			Tokens: []*scanner.Token{
				{Type: tokentype.Var, Lexeme: `var`},
				{Type: tokentype.Identifier, Lexeme: `language`},
				{Type: tokentype.Equal, Lexeme: `=`},
				{Type: tokentype.String, Lexeme: `"lox"`, Literal: literal.String([]runes.Rune("lox"))},
				{Type: tokentype.Semicolon, Lexeme: `;`},
				{Type: tokentype.EOF},
			},
		},
		{
			Name:   "error/unterminated string",
			Source: `"Hello!`,
			Tokens: nil,
			Err:    scanner.ErrUnterminatedString,
		},
		{
			Name:   "success/assign number to variable",
			Source: `var pi = 873.32;`,
			Tokens: []*scanner.Token{
				{Type: tokentype.Var, Lexeme: `var`},
				{Type: tokentype.Identifier, Lexeme: `pi`},
				{Type: tokentype.Equal, Lexeme: `=`},
				{Type: tokentype.Number, Lexeme: `873.32`, Literal: literal.Number(873.32)},
				{Type: tokentype.Semicolon, Lexeme: `;`},
				{Type: tokentype.EOF},
			},
		},
		{
			Name: "success/multiline string",
			Source: `"First line
			Second line"`,
			Tokens: []*scanner.Token{
				{Type: tokentype.String, Lexeme: `"First line
			Second line"`, Literal: literal.String([]runes.Rune(`First line
			Second line`)), Line: 1},
				{Type: tokentype.EOF, Line: 1},
			},
		},
		{
			Name:   "success/single character punctuation",
			Source: `(){}=;*+-/.,!`,
			Tokens: []*scanner.Token{
				{Type: tokentype.LeftParen, Lexeme: `(`},
				{Type: tokentype.RightParen, Lexeme: `)`},
				{Type: tokentype.LeftBrace, Lexeme: `{`},
				{Type: tokentype.RightBrace, Lexeme: `}`},
				{Type: tokentype.Equal, Lexeme: `=`},
				{Type: tokentype.Semicolon, Lexeme: `;`},
				{Type: tokentype.Star, Lexeme: `*`},
				{Type: tokentype.Plus, Lexeme: `+`},
				{Type: tokentype.Minus, Lexeme: `-`},
				{Type: tokentype.Slash, Lexeme: `/`},
				{Type: tokentype.Dot, Lexeme: `.`},
				{Type: tokentype.Comma, Lexeme: `,`},
				{Type: tokentype.Bang, Lexeme: `!`},
				{Type: tokentype.EOF},
			},
		},
		{
			Name:   "success/double character punctuation",
			Source: `= <= >= == !=`,
			Tokens: []*scanner.Token{
				{Type: tokentype.Equal, Lexeme: `=`},
				{Type: tokentype.LessEqual, Lexeme: `<=`},
				{Type: tokentype.GreaterEqual, Lexeme: `>=`},
				{Type: tokentype.EqualEqual, Lexeme: `==`},
				{Type: tokentype.BangEqual, Lexeme: `!=`},
				{Type: tokentype.EOF},
			},
		},
		{
			Name: "success/multiline source",
			Source: `
				// add_two adds returns the number incremented by 2.
				fun add_two(n) {
					return n + 2;
				}
			`,
			Tokens: []*scanner.Token{
				{Type: tokentype.Fun, Lexeme: `fun`, Line: 2},
				{Type: tokentype.Identifier, Lexeme: `add_two`, Line: 2},
				{Type: tokentype.LeftParen, Lexeme: `(`, Line: 2},
				{Type: tokentype.Identifier, Lexeme: `n`, Line: 2},
				{Type: tokentype.RightParen, Lexeme: `)`, Line: 2},
				{Type: tokentype.LeftBrace, Lexeme: `{`, Line: 2},
				{Type: tokentype.Return, Lexeme: `return`, Line: 3},
				{Type: tokentype.Identifier, Lexeme: `n`, Line: 3},
				{Type: tokentype.Plus, Lexeme: `+`, Line: 3},
				{Type: tokentype.Number, Lexeme: `2`, Literal: literal.Number(2), Line: 3},
				{Type: tokentype.Semicolon, Lexeme: `;`, Line: 3},
				{Type: tokentype.RightBrace, Lexeme: `}`, Line: 4},
				{Type: tokentype.EOF, Line: 5},
			},
		},
		{
			Name:   "success/number literal",
			Source: `4.`,
			Tokens: []*scanner.Token{
				{Type: tokentype.Number, Lexeme: `4`, Literal: literal.Number(4)},
				{Type: tokentype.Dot, Lexeme: `.`},
				{Type: tokentype.EOF},
			},
		},
		{
			Name:   "error/unexpected rune",
			Source: `'`,
			Err:    scanner.ErrUnexpectedRune,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			tokens, err := scanner.New(test.Source).ScanTokens()
			if test.Err != nil {
				require.ErrorIs(t, err, test.Err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.Tokens, tokens)
			}
		})
	}
}
