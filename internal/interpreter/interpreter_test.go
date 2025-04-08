package interpreter_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/interpreter"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
)

func TestInterpreter_Evaluate(t *testing.T) {
	t.Parallel()

	type Test struct {
		Name          string
		Source        string
		ExpectedValue loxtype.Type
		ExpectedError error
	}

	tests := []Test{
		{
			Name:          "success/negate_integer",
			Source:        `- 4;`,
			ExpectedValue: loxtype.Number(-4),
		},
		{
			Name:          "error/negate_string",
			Source:        `-"Hello";`,
			ExpectedError: interpreter.ErrNonNumericType,
		},
		{
			Name:          "success/truthiness/nil_is_falsey",
			Source:        `!nil;`,
			ExpectedValue: loxtype.Boolean(true),
		},
		{
			Name: "success/assignment",
			Source: `
			var a = "hello";
			a = "world";
			a;
			`,
			ExpectedValue: loxtype.String("world"),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			i := &interpreter.Interpreter{}
			tokens, err := scanner.New(test.Source).ScanTokens()
			require.NoError(t, err)
			stmts, err := parser.New(tokens).Parse()
			require.NoError(t, err)

			env := environment.New()
			err = i.Interpret(env, stmts)

			if test.ExpectedError != nil {
				require.ErrorIs(t, err, test.ExpectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}
