package literal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/literal"
	"github.com/matt-hoiland/glox/internal/scanner/runes"
)

func TestBoolean_String(t *testing.T) {
	bt, bf := literal.Boolean(true), literal.Boolean(false)
	assert.Equal(t, "true", bt.String())
	assert.Equal(t, "false", bf.String())
}

func TestNil_String(t *testing.T) {
	assert.Equal(t, "nil", literal.Nil{}.String())
}

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

func TestString_String(t *testing.T) {
	t.Parallel()

	stdString := "Hello, world!"
	myString := literal.String(stdString)
	value := myString.String()
	assert.Equal(t, stdString, value)
}
