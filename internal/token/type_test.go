package token_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/token"
)

func TestType_String(t *testing.T) {
	t.Parallel()

	t.Run("in range", func(t *testing.T) {
		t.Parallel()
		tt := token.TypeBangEqual
		s := tt.String()
		assert.Equal(t, "TypeBangEqual", s)
	})

	t.Run("out of range", func(t *testing.T) {
		t.Parallel()
		tt := token.Type(-420)
		s := tt.String()
		assert.Equal(t, "Type(-420)", s)
	})
}
