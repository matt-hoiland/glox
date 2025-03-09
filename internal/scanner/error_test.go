package scanner_test

import (
	"testing"

	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	t.Parallel()

	var err error = &scanner.Error{
		Line:  42,
		Where: "blah",
		Err:   scanner.ErrUnterminatedString,
	}

	s := err.Error()
	require.ErrorIs(t, err, scanner.ErrUnterminatedString)
	assert.Equal(t, "[line 42] Errorblah: unterminated string", s)
}
