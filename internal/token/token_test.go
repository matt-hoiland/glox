package token_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/token"
	"github.com/matt-hoiland/glox/internal/token/tokentype"
)

func TestToken_String(t *testing.T) {
	t.Parallel()

	t.Run("with literal", func(t *testing.T) {
		t.Parallel()

		token := &token.Token{
			Type:    tokentype.String,
			Lexeme:  `"Hello, world!"`,
			Literal: literal.String("Hello, world!"),
		}

		assert.Equal(t, `String "Hello, world!" 'Hello, world!'`, token.String())
	})

	t.Run("without literal", func(t *testing.T) {
		t.Parallel()

		token := &token.Token{
			Type:   tokentype.And,
			Lexeme: `and`,
		}

		assert.Equal(t, `And and`, token.String())
	})
}
