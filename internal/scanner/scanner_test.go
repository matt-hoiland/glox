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
			Name:   "success/empty_source",
			Source: ``,
			Tokens: []*token.Token{
				{Type: token.TypeEOF, Line: 1},
			},
		},
		{
			Name:   "success/assign_string_to_variable",
			Source: `var language = "lox";`,
			Tokens: []*token.Token{
				{Type: token.TypeVar, Lexeme: `var`, Line: 1},
				{Type: token.TypeIdentifier, Lexeme: `language`, Line: 1},
				{Type: token.TypeEqual, Lexeme: `=`, Line: 1},
				{Type: token.TypeString, Lexeme: `"lox"`, Literal: loxtype.String([]runes.Rune("lox")), Line: 1},
				{Type: token.TypeSemicolon, Lexeme: `;`, Line: 1},
				{Type: token.TypeEOF, Line: 1},
			},
		},
		{
			Name:   "error/unterminated_string",
			Source: `"Hello!`,
			Tokens: nil,
			Err:    scanner.ErrUnterminatedString,
		},
		{
			Name:   "success/assign_number_to_variable",
			Source: `var pi = 873.32;`,
			Tokens: []*token.Token{
				{Type: token.TypeVar, Lexeme: `var`, Line: 1},
				{Type: token.TypeIdentifier, Lexeme: `pi`, Line: 1},
				{Type: token.TypeEqual, Lexeme: `=`, Line: 1},
				{Type: token.TypeNumber, Lexeme: `873.32`, Literal: loxtype.Number(873.32), Line: 1},
				{Type: token.TypeSemicolon, Lexeme: `;`, Line: 1},
				{Type: token.TypeEOF, Line: 1},
			},
		},
		{
			Name: "success/multiline_string",
			Source: `"First line
			Second line"`,
			Tokens: []*token.Token{
				{
					Type: token.TypeString, Lexeme: `"First line
			Second line"`,
					Literal: loxtype.String([]runes.Rune(`First line
			Second line`)),
					Line: 2, // The line number corresponds to the last character.
				},
				{Type: token.TypeEOF, Line: 2},
			},
		},
		{
			Name:   "success/single_character_punctuation",
			Source: `(){}=;*+-/.,!`,
			Tokens: []*token.Token{
				{Type: token.TypeLeftParen, Lexeme: `(`, Line: 1},
				{Type: token.TypeRightParen, Lexeme: `)`, Line: 1},
				{Type: token.TypeLeftBrace, Lexeme: `{`, Line: 1},
				{Type: token.TypeRightBrace, Lexeme: `}`, Line: 1},
				{Type: token.TypeEqual, Lexeme: `=`, Line: 1},
				{Type: token.TypeSemicolon, Lexeme: `;`, Line: 1},
				{Type: token.TypeStar, Lexeme: `*`, Line: 1},
				{Type: token.TypePlus, Lexeme: `+`, Line: 1},
				{Type: token.TypeMinus, Lexeme: `-`, Line: 1},
				{Type: token.TypeSlash, Lexeme: `/`, Line: 1},
				{Type: token.TypeDot, Lexeme: `.`, Line: 1},
				{Type: token.TypeComma, Lexeme: `,`, Line: 1},
				{Type: token.TypeBang, Lexeme: `!`, Line: 1},
				{Type: token.TypeEOF, Line: 1},
			},
		},
		{
			Name:   "success/double_character_punctuation",
			Source: `= <= >= == !=`,
			Tokens: []*token.Token{
				{Type: token.TypeEqual, Lexeme: `=`, Line: 1},
				{Type: token.TypeLessEqual, Lexeme: `<=`, Line: 1},
				{Type: token.TypeGreaterEqual, Lexeme: `>=`, Line: 1},
				{Type: token.TypeEqualEqual, Lexeme: `==`, Line: 1},
				{Type: token.TypeBangEqual, Lexeme: `!=`, Line: 1},
				{Type: token.TypeEOF, Line: 1},
			},
		},
		{
			Name: "success/multiline_source",
			Source: `
				// add_two adds returns the number incremented by 2.
				fun add_two(n) {
					return n + 2;
				}
			`,
			Tokens: []*token.Token{
				{Type: token.TypeFun, Lexeme: `fun`, Line: 3},
				{Type: token.TypeIdentifier, Lexeme: `add_two`, Line: 3},
				{Type: token.TypeLeftParen, Lexeme: `(`, Line: 3},
				{Type: token.TypeIdentifier, Lexeme: `n`, Line: 3},
				{Type: token.TypeRightParen, Lexeme: `)`, Line: 3},
				{Type: token.TypeLeftBrace, Lexeme: `{`, Line: 3},
				{Type: token.TypeReturn, Lexeme: `return`, Line: 4},
				{Type: token.TypeIdentifier, Lexeme: `n`, Line: 4},
				{Type: token.TypePlus, Lexeme: `+`, Line: 4},
				{Type: token.TypeNumber, Lexeme: `2`, Literal: loxtype.Number(2), Line: 4},
				{Type: token.TypeSemicolon, Lexeme: `;`, Line: 4},
				{Type: token.TypeRightBrace, Lexeme: `}`, Line: 5},
				{Type: token.TypeEOF, Line: 6},
			},
		},
		{
			Name:   "success/number_literal",
			Source: `4.`,
			Tokens: []*token.Token{
				{Type: token.TypeNumber, Lexeme: `4`, Literal: loxtype.Number(4), Line: 1},
				{Type: token.TypeDot, Lexeme: `.`, Line: 1},
				{Type: token.TypeEOF, Line: 1},
			},
		},
		{
			Name:   "error/unexpected_rune",
			Source: `'`,
			Err:    scanner.ErrUnexpectedRune,
		},
		{
			Name:   "success/nil_literal",
			Source: `nil`,
			Tokens: []*token.Token{
				{Type: token.TypeNil, Lexeme: `nil`, Line: 1},
				{Type: token.TypeEOF, Line: 1},
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
