package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/errors"
)

func TestError_Error(t *testing.T) {
	t.Parallel()

	var err error = &errors.Error{
		Line:  42,
		Where: "blah",
		Err:   assert.AnError,
	}

	s := err.Error()
	require.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, "[line 42] Errorblah: assert.AnError general error for testing", s)
}
