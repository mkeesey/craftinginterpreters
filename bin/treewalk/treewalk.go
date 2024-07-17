package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mkeesey/craftinginterpreters/pkg/ast"
	"github.com/mkeesey/craftinginterpreters/pkg/failure"
	"github.com/mkeesey/craftinginterpreters/pkg/parser"
	"github.com/mkeesey/craftinginterpreters/pkg/scanner"
)

var (
	reporter = &failure.Reporter{}
	visitor  = ast.NewInterpreter(reporter)
)

func main() {
	var err error
	if len(os.Args) == 2 {
		err = runFile(os.Args[1])
	} else if len(os.Args) == 1 {
		err = runPrompt()
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s [script]\n", os.Args[0])
		os.Exit(64)
	}

	if err != nil {
		var outputter ErrorOutputter
		if errors.As(err, &outputter) {
			outputter.Output(os.Stderr)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		var compileError *CompileError
		var runtimeError *RuntimeError
		if errors.As(err, &compileError) {
			os.Exit(65)
		} else if errors.As(err, &runtimeError) {
			os.Exit(70)
		} else {
			os.Exit(1)
		}
	}
}

func runFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	return run(file)
}

func runPrompt() error {
	inputReader := bufio.NewReader(os.Stdin)
	reader := strings.NewReader("")
	for {
		fmt.Printf("> ")
		line, err := inputReader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil // done
			}
			return fmt.Errorf("failed to read line: %w", err)
		}

		reader.Reset(line)
		err = run(reader)
		if err != nil {
			var outputter ErrorOutputter
			if errors.As(err, &outputter) {
				outputter.Output(os.Stderr)
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		}
	}
}

func run(reader io.Reader) error {

	scan := scanner.NewScanner(reader, reporter)
	tokens := scan.ScanTokens()
	if reporter.HasFailed() {
		return NewCompileError("")
	}

	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()
	// TODO - replace chain of errs with reporter usage
	if err != nil {
		return NewCompileError(err.Error())
	}

	resolver := ast.NewResolver(visitor, reporter)
	resolver.Resolve(statements)
	if reporter.HasFailed() {
		return NewCompileError("")
	}

	visitor.Interpret(statements)
	if reporter.HasFailed() {
		return NewRuntimeError()
	}
	return nil
}

type ErrorOutputter interface {
	Output(w io.Writer)
}

type CompileError struct {
	Message string
}

func NewCompileError(message string) error {
	return &CompileError{Message: message}
}

func (c *CompileError) Error() string {
	return fmt.Sprintf("CompileError %s", c.Message)
}

func (c *CompileError) Output(w io.Writer) {
	if c.Message != "" {
		fmt.Fprintf(w, "%s\n", c.Message)
	}
}

type RuntimeError struct {
}

func NewRuntimeError() error {
	return &RuntimeError{}
}

func (r *RuntimeError) Error() string {
	return "RuntimeError"
}

func (r *RuntimeError) Output(w io.Writer) {
}
