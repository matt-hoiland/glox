package scanner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/runes"
	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/matt-hoiland/glox/internal/token"
)

func TestScanner_ScanTokens(t *testing.T) {
	t.Parallel()

	type Test struct {
		Name   string
		Source string
		Tokens []*token.Token
		Err    error
	}

	tests := []Test{
		{
			Name:   "success/empty source",
			Source: ``,
			Tokens: []*token.Token{
				{Type: token.TypeEOF},
			},
		},
		{
			Name:   "success/assign string to variable",
			Source: `var language = "lox";`,
			Tokens: []*token.Token{
				{Type: token.TypeVar, Lexeme: `var`},
				{Type: token.TypeIdentifier, Lexeme: `language`},
				{Type: token.TypeEqual, Lexeme: `=`},
				{Type: token.TypeString, Lexeme: `"lox"`, Literal: loxtype.String([]runes.Rune("lox"))},
				{Type: token.TypeSemicolon, Lexeme: `;`},
				{Type: token.TypeEOF},
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
			Tokens: []*token.Token{
				{Type: token.TypeVar, Lexeme: `var`},
				{Type: token.TypeIdentifier, Lexeme: `pi`},
				{Type: token.TypeEqual, Lexeme: `=`},
				{Type: token.TypeNumber, Lexeme: `873.32`, Literal: loxtype.Number(873.32)},
				{Type: token.TypeSemicolon, Lexeme: `;`},
				{Type: token.TypeEOF},
			},
		},
		{
			Name: "success/multiline string",
			Source: `"First line
			Second line"`,
			Tokens: []*token.Token{
				{Type: token.TypeString, Lexeme: `"First line
			Second line"`, Literal: loxtype.String([]runes.Rune(`First line
			Second line`)), Line: 1},
				{Type: token.TypeEOF, Line: 1},
			},
		},
		{
			Name:   "success/single character punctuation",
			Source: `(){}=;*+-/.,!`,
			Tokens: []*token.Token{
				{Type: token.TypeLeftParen, Lexeme: `(`},
				{Type: token.TypeRightParen, Lexeme: `)`},
				{Type: token.TypeLeftBrace, Lexeme: `{`},
				{Type: token.TypeRightBrace, Lexeme: `}`},
				{Type: token.TypeEqual, Lexeme: `=`},
				{Type: token.TypeSemicolon, Lexeme: `;`},
				{Type: token.TypeStar, Lexeme: `*`},
				{Type: token.TypePlus, Lexeme: `+`},
				{Type: token.TypeMinus, Lexeme: `-`},
				{Type: token.TypeSlash, Lexeme: `/`},
				{Type: token.TypeDot, Lexeme: `.`},
				{Type: token.TypeComma, Lexeme: `,`},
				{Type: token.TypeBang, Lexeme: `!`},
				{Type: token.TypeEOF},
			},
		},
		{
			Name:   "success/double character punctuation",
			Source: `= <= >= == !=`,
			Tokens: []*token.Token{
				{Type: token.TypeEqual, Lexeme: `=`},
				{Type: token.TypeLessEqual, Lexeme: `<=`},
				{Type: token.TypeGreaterEqual, Lexeme: `>=`},
				{Type: token.TypeEqualEqual, Lexeme: `==`},
				{Type: token.TypeBangEqual, Lexeme: `!=`},
				{Type: token.TypeEOF},
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
			Tokens: []*token.Token{
				{Type: token.TypeFun, Lexeme: `fun`, Line: 2},
				{Type: token.TypeIdentifier, Lexeme: `add_two`, Line: 2},
				{Type: token.TypeLeftParen, Lexeme: `(`, Line: 2},
				{Type: token.TypeIdentifier, Lexeme: `n`, Line: 2},
				{Type: token.TypeRightParen, Lexeme: `)`, Line: 2},
				{Type: token.TypeLeftBrace, Lexeme: `{`, Line: 2},
				{Type: token.TypeReturn, Lexeme: `return`, Line: 3},
				{Type: token.TypeIdentifier, Lexeme: `n`, Line: 3},
				{Type: token.TypePlus, Lexeme: `+`, Line: 3},
				{Type: token.TypeNumber, Lexeme: `2`, Literal: loxtype.Number(2), Line: 3},
				{Type: token.TypeSemicolon, Lexeme: `;`, Line: 3},
				{Type: token.TypeRightBrace, Lexeme: `}`, Line: 4},
				{Type: token.TypeEOF, Line: 5},
			},
		},
		{
			Name:   "success/number literal",
			Source: `4.`,
			Tokens: []*token.Token{
				{Type: token.TypeNumber, Lexeme: `4`, Literal: loxtype.Number(4)},
				{Type: token.TypeDot, Lexeme: `.`},
				{Type: token.TypeEOF},
			},
		},
		{
			Name:   "error/unexpected rune",
			Source: `'`,
			Err:    scanner.ErrUnexpectedRune,
		},
		{
			Name:   "success/nil literal",
			Source: `nil`,
			Tokens: []*token.Token{
				{Type: token.TypeNil, Lexeme: `nil`},
				{Type: token.TypeEOF},
			},
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
