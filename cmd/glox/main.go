package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/matt-hoiland/glox/internal/constants/exit"
	"github.com/matt-hoiland/glox/internal/scanner"
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
	scanner := scanner.New(code)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}
	for _, token := range tokens {
		fmt.Println(token.String())
	}
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
