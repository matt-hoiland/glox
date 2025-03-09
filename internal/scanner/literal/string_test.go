package literal_test

import (
	"testing"

	"github.com/matt-hoiland/glox/internal/scanner/literal"
	"github.com/stretchr/testify/assert"
)

func TestString_String(t *testing.T) {
	t.Parallel()

	stdString := "Hello, world!"
	myString := literal.String(stdString)
	value := myString.String()
	assert.Equal(t, stdString, value)
}
