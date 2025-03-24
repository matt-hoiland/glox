package token_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/token"
)

func TestToken_String(t *testing.T) {
	t.Parallel()

	t.Run("with literal", func(t *testing.T) {
		t.Parallel()

		token := &token.Token{
			Type:    token.TypeString,
			Lexeme:  `"Hello, world!"`,
			Literal: literal.String("Hello, world!"),
		}

		assert.Equal(t, `TypeString "Hello, world!" 'Hello, world!'`, token.String())
	})

	t.Run("without literal", func(t *testing.T) {
		t.Parallel()

		token := &token.Token{
			Type:   token.TypeAnd,
			Lexeme: `and`,
		}

		assert.Equal(t, `TypeAnd and`, token.String())
	})
}
