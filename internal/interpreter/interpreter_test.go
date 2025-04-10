package interpreter_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/interpreter"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
)

func TestInterpreter_Evaluate(t *testing.T) {
	// t.Parallel()

	type Test struct {
		Name           string
		Source         string
		ExpectedOutput string
		ExpectedError  error
	}

	tests := []Test{
		{
			Name:           "success/negate_integer",
			Source:         `print (- 4);`,
			ExpectedOutput: "-4\n",
		},
		{
			Name:          "error/negate_string",
			Source:        `print (-"Hello");`,
			ExpectedError: interpreter.ErrNonNumericType,
		},
		{
			Name:           "success/truthiness/nil_is_falsey",
			Source:         `print (!nil);`,
			ExpectedOutput: "true\n",
		},
		{
			Name: "success/assignment",
			Source: `
			var a = "hello";
			a = "world";
			print a;
			`,
			ExpectedOutput: "world\n",
		},
		{
			Name: "success/lexical_scope",
			Source: `
				var a = 4;
				var shadowed = "shadow";
				{
					a = 6;
					var shadowed = "block";
					print shadowed;
				}
				print a;
				print shadowed;
			`,
			ExpectedOutput: stripIndentation(`
				block
				6
				shadow
			`),
		},
		{
			Name: "success/nystrom/lexical_scope",
			Source: `
				var a = "global a";
				var b = "global b";
				var c = "global c";
				{
					var a = "outer a";
					var b = "outer b";
					{
						var a = "inner a";
						print a;
						print b;
						print c;
					}
					print a;
					print b;
					print c;
				}
				print a;
				print b;
				print c;
			`,
			ExpectedOutput: stripIndentation(`
				inner a
				outer b
				global c
				outer a
				outer b
				global c
				global a
				global b
				global c
			`),
		},
		{
			Name: "success/if_statements",
			Source: `
				if (true) {
					print "Hello";
				}

				if (false) {
					print "Matt!";
				} else {
					print "World!";
				}
			`,
			ExpectedOutput: stripIndentation(`
				Hello
				World!
			`),
		},
		{
			Name: "success/short_circuiting",
			Source: `
				print "hi" or 2;     // "hi".
				print nil or "yes";  // "yes".
				print nil and "bye"; // "nil".
			`,
			ExpectedOutput: stripIndentation(`
				hi
				yes
				nil
			`),
		},
		{
			Name: "success/simple_while_loop",
			Source: `
				var i = 0;
				while (i < 5) {
					print "Hello";
					i = i + 1;
				}
			`,
			ExpectedOutput: stripIndentation(`
				Hello
				Hello
				Hello
				Hello
				Hello
			`),
		},
		{
			Name: "success/for_loop/fibonacci",
			Source: `
				var a = 0;
				var temp;

				for (var b = 1; a < 10000; b = temp + b) {
					print a;
					temp = a;
					a = b;
				}
			`,
			ExpectedOutput: stripIndentation(`
				0
				1
				1
				2
				3
				5
				8
				13
				21
				34
				55
				89
				144
				233
				377
				610
				987
				1597
				2584
				4181
				6765
			`),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// t.Parallel()

			var bob strings.Builder
			i := interpreter.New(&bob)
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
			assert.Equal(t, test.ExpectedOutput, bob.String())
		})
	}
}

func stripIndentation(s string) string {
	var (
		bob             strings.Builder
		passedFirstLine bool
	)
	for line := range strings.Lines(s) {
		if !passedFirstLine {
			passedFirstLine = true
			continue
		}
		bob.WriteString(strings.TrimLeft(line, "\t "))
	}
	return bob.String()
}
