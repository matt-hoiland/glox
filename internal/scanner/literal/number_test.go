package literal_test

import (
	"testing"

	"github.com/matt-hoiland/glox/internal/scanner/literal"
	"github.com/matt-hoiland/glox/internal/scanner/runes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNumber(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		s := []runes.Rune("3.14")
		n, err := literal.ParseNumber(s)
		require.NoError(t, err)
		assert.Equal(t, literal.Number(3.14), n)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		s := []runes.Rune("banana")
		n, err := literal.ParseNumber(s)
		require.Error(t, err)
		assert.Zero(t, n)
	})
}

func TestNumber_String(t *testing.T) {
	n := literal.Number(3.14)
	s := n.String()
	assert.Equal(t, "3.14", s)
}
