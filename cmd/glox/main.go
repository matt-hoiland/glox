package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/matt-hoiland/glox/internal/ast"
	"github.com/matt-hoiland/glox/internal/constants/exit"
	"github.com/matt-hoiland/glox/internal/interpreter"
	"github.com/matt-hoiland/glox/internal/loxtype"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/matt-hoiland/glox/internal/token"
)

func main() {
	args := os.Args[1:]

	var err error
	switch {
	case len(args) > 1:
		fmt.Fprintln(os.Stdout, "Usage: glox [script]")
		os.Exit(exit.Usage)
	case len(args) == 1:
		err = runFile(args[0])
	default:
		runPrompt()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(exit.DataErr)
	}
}

func runFile(filename string) error {
	var (
		data []byte
		err  error
	)
	if data, err = os.ReadFile(filename); err != nil {
		return fmt.Errorf("could not read file '%s': %w", filename, err)
	}
	return interpreter.New(os.Stdout).Run(string(data))
}

func runPrompt() {
	var (
		w          io.Writer = os.Stdout
		i                    = interpreter.New(w)
		reader               = bufio.NewScanner(os.Stdin)
		lineNumber           = 0
		err        error
	)
	for {
		lineNumber++
		fmt.Fprintf(os.Stdin, "#%3d > ", lineNumber)
		if !reader.Scan() {
			break
		}
		line := reader.Text()

		var (
			tokens []*token.Token
			s      = scanner.New(line, scanner.WithStartingLine(lineNumber))
		)
		if tokens, err = s.ScanTokens(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		var (
			stmts []ast.Stmt
			p     = parser.New(tokens, parser.InREPLMode())
		)
		if stmts, err = p.Parse(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		if len(stmts) == 1 {
			if exprStmt, ok := stmts[0].(*ast.ExpressionStmt); ok {
				var value loxtype.Type
				if value, err = i.Evaluate(exprStmt.Expression); err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					continue
				}
				fmt.Fprintln(os.Stdout, value.String())
				continue
			}
		}

		if err = i.Interpret(stmts); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
