package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/matt-hoiland/glox/internal/constants/exit"
	"github.com/matt-hoiland/glox/internal/interpreter"
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

func runFile(filename string) (err error) {
	var data []byte
	if data, err = os.ReadFile(filename); err != nil {
		return fmt.Errorf("could not read file '%s': %w", filename, err)
	}
	return interpreter.New().Run(string(data))
}

func runPrompt() error {
	var (
		i      = interpreter.New()
		reader = bufio.NewScanner(os.Stdin)
	)
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		line := reader.Text()
		if err := i.Run(line); err != nil {
			return err
		}
	}
	return nil
}
