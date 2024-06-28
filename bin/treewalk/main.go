package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mkeesey/craftinginterpreters/pkg/scanner"
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
	scan := scanner.NewScanner(reader)
	tokens, err := scan.ScanTokens()
	if err != nil {
		return fmt.Errorf("Error scanning tokens: %w", err)
	}
	fmt.Printf("%+v\n", tokens)
	return nil
}
