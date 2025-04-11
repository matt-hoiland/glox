package interpreter_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matt-hoiland/glox/internal/interpreter"
)

func TestInterpreter_Evaluate(t *testing.T) {
	t.Parallel()

	type Test struct {
		Name          string
		Source        string
		Output        string
		ExpectedError error
	}

	tests := []Test{
		{
			Name:   "success/negate_integer",
			Source: `print (- 4);`,
			Output: "-4\n",
		},
		{
			Name:          "error/negate_string",
			Source:        `print (-"Hello");`,
			ExpectedError: interpreter.ErrNonNumericType,
		},
		{
			Name:   "success/truthiness/nil_is_falsey",
			Source: `print (!nil);`,
			Output: "true\n",
		},
		{
			Name: "success/assignment",
			Source: `
			var a = "hello";
			a = "world";
			print a;
			`,
			Output: "world\n",
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
			Output: dedent(`
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
			Output: dedent(`
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
			Output: dedent(`
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
			Output: dedent(`
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
			Output: dedent(`
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
			Output: dedent(`
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
		{
			Name: "success/function_call/clock",
			Source: `
				var mils = clock();
				if (mils > 0) {
					print "success";
				}
			`,
			Output: dedent(`
				success
			`),
		},
		{
			Name: "success/functions/count_to_three",
			Source: `
				fun count(n) {
					if (n > 1) count(n - 1);
					print n;
				}

				count(3);
			`,
			Output: dedent(`
				1
				2
				3
			`),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			var bob strings.Builder
			err := interpreter.New(&bob).Run(test.Source)

			if test.ExpectedError != nil {
				require.ErrorIs(t, err, test.ExpectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.Output, bob.String())
		})
	}
}

func dedent(s string) string {
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
