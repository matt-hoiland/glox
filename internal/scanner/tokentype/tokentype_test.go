package tokentype_test

import (
	"testing"

	"github.com/matt-hoiland/glox/internal/scanner/tokentype"
	"github.com/stretchr/testify/assert"
)

func TestTokenType_String(t *testing.T) {
	t.Parallel()

	t.Run("in range", func(t *testing.T) {
		t.Parallel()
		tt := tokentype.BangEqual
		s := tt.String()
		assert.Equal(t, "BangEqual", s)
	})

	t.Run("out of range", func(t *testing.T) {
		t.Parallel()
		tt := tokentype.TokenType(-420)
		s := tt.String()
		assert.Equal(t, "TokenType(-420)", s)
	})
}
