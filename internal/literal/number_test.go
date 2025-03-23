package literal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/scanner/runes"
)

func TestParseNumber(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		s := []runes.Rune("3.14")
		n := literal.ParseNumber(s)
		assert.Equal(t, literal.Number(3.14), n)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		s := []runes.Rune("banana")
		n := literal.ParseNumber(s)
		assert.Zero(t, n)
	})
}

func TestNumber_String(t *testing.T) {
	n := literal.Number(3.14)
	s := n.String()
	assert.Equal(t, "3.14", s)
}
