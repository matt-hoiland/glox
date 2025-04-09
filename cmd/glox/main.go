package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/matt-hoiland/glox/internal/constants/exit"
	"github.com/matt-hoiland/glox/internal/environment"
	"github.com/matt-hoiland/glox/internal/interpreter"
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
		err = runPrompt()
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
	return interpreter.New(os.Stdout).Run(environment.New(), string(data))
}

func runPrompt() error {
	var (
		w         io.Writer = os.Stdout
		i                   = interpreter.New(w)
		reader              = bufio.NewScanner(os.Stdin)
		globalEnv           = environment.New()
	)
	for {
		fmt.Fprint(os.Stdin, "> ")
		if !reader.Scan() {
			break
		}
		line := reader.Text()
		if err := i.Run(globalEnv, line); err != nil {
			return err
		}
	}
	return nil
}
