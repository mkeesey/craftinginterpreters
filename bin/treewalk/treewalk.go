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
	visitor = ast.NewInterpreter()
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
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func run(reader io.Reader) error {
	reporter := &failure.Reporter{}

	scan := scanner.NewScanner(reader, reporter)
	tokens := scan.ScanTokens()
	if reporter.HasFailed() {
		return errors.New("Scanner failed")
	}

	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()
	// TODO - replace chain of errs with reporter usage
	if err != nil {
		return err
	}

	resolver := ast.NewResolver(visitor, reporter)
	resolver.Resolve(statements)
	if reporter.HasFailed() {
		return errors.New("Resolver failed")
	}

	err = visitor.Interpret(statements)
	if err != nil {
		// TODO - distinguish runtime from parse errors for exit codes
		return err
	}
	return nil
}
