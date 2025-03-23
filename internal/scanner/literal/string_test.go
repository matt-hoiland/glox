package literal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matt-hoiland/glox/internal/scanner/literal"
)

func TestString_String(t *testing.T) {
	t.Parallel()

	stdString := "Hello, world!"
	myString := literal.String(stdString)
	value := myString.String()
	assert.Equal(t, stdString, value)
}
