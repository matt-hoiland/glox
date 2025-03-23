package runes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/scanner/runes"
)

func TestRune_IsAlpha(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()

		text := []runes.Rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
		for _, r := range text {
			assert.True(t, r.IsAlpha())
		}
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()
		text := []runes.Rune("12344567890-*`+=\n")
		for _, r := range text {
			assert.False(t, r.IsAlpha())
		}
	})
}

func TestRune_IsAlphaNumeric(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()
		text := []runes.Rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789")
		for _, r := range text {
			assert.True(t, r.IsAlphaNumeric())
		}
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()
		text := []runes.Rune(".$#*-+")
		for _, r := range text {
			assert.False(t, r.IsAlphaNumeric())
		}
	})
}

func TestRune_IsDigit(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()
		text := []runes.Rune("0123456789")
		for _, r := range text {
			assert.True(t, r.IsDigit())
		}
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()
		text := []runes.Rune("abcdefghijklmnopqrstuvwxyz.$#*-+_")
		for _, r := range text {
			assert.False(t, r.IsDigit())
		}
	})
}
