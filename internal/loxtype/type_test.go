package loxtype_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/runes"
)

func TestBoolean_String(t *testing.T) {
	bt, bf := loxtype.Boolean(true), loxtype.Boolean(false)
	assert.Equal(t, "true", bt.String())
	assert.Equal(t, "false", bf.String())
}

func TestNil_String(t *testing.T) {
	assert.Equal(t, "nil", loxtype.Nil{}.String())
}

func TestParseNumber(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		s := []runes.Rune("3.14")
		n := loxtype.ParseNumber(s)
		assert.InEpsilon(t, 3.14, float64(n), 0.001)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		s := []runes.Rune("banana")
		n := loxtype.ParseNumber(s)
		assert.Zero(t, n)
	})
}

func TestNumber_String(t *testing.T) {
	n := loxtype.Number(3.14)
	s := n.String()
	assert.Equal(t, "3.14", s)
}

func TestString_String(t *testing.T) {
	t.Parallel()

	stdString := "Hello, world!"
	myString := loxtype.String(stdString)
	value := myString.String()
	assert.Equal(t, stdString, value)
}
