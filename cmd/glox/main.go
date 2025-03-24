package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/matt-hoiland/glox/internal/astprinter"
	"github.com/matt-hoiland/glox/internal/constants/exit"
	"github.com/matt-hoiland/glox/internal/expr"
	"github.com/matt-hoiland/glox/internal/parser"
	"github.com/matt-hoiland/glox/internal/scanner"
	"github.com/matt-hoiland/glox/internal/token"
)

func main() {
	args := os.Args[1:]

	var err error
	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(exit.Usage)
	} else if len(args) == 1 {
		err = runFile(args[0])
	} else {
		err = runPrompt()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(exit.DataErr)
	}
}

func run(code string) error {
	var (
		tokens []*token.Token
		ast    expr.Expr
		err    error
	)

	if tokens, err = scanner.New(code).ScanTokens(); err != nil {
		return err
	}

	if ast, err = parser.New(tokens).Parse(); err != nil {
		return err
	}

	fmt.Println(astprinter.Print(ast))
	return nil
}

func runFile(filename string) (err error) {
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		return fmt.Errorf("could not read file '%s': %w", filename, err)
	}
	return run(string(data))
}

func runPrompt() error {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		line := reader.Text()
		if err := run(line); err != nil {
			return err
		}
	}
	return nil
}
