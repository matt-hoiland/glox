package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/matt-hoiland/glox/internal/constants/exit"
	"github.com/matt-hoiland/glox/internal/scanner"
)

const ExitUsage = 64

var hadError bool

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(exit.Usage)
	} else if len(args) == 1 {
		if err := runFile(args[0]); err != nil {
			Error(0, err.Error())
			os.Exit(exit.DataErr)
		}
	} else {
		runPrompt()
	}
}

func run(code string) {
	scanner := scanner.New(code)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	for _, token := range tokens {
		fmt.Println(token.String())
	}
}

func runFile(filename string) (err error) {
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		return fmt.Errorf("could not read file '%s': %w", filename, err)
	}
	run(string(data))

	// Indicate an error in the exit code.
	if hadError {
		os.Exit(exit.DataErr)
	}
	return nil
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		line := reader.Text()
		run(line)
	}
}

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	hadError = true
}
